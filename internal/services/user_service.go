package services

import (
	"context"
	"time"

	"laundry-backend/internal/dto"
	"laundry-backend/internal/models"
	"laundry-backend/internal/repositories"
	"laundry-backend/pkg/response"
	"laundry-backend/pkg/utils"
)

// UserService defines the contract for business logic related to users.
type UserService interface {
	RegisterUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserDetailResponse, error)
	GetUsers(ctx context.Context, page, perPage int, search, role string, status int) (*dto.UserListResponse, error)
	GetUserProfile(ctx context.Context, id int64) (*dto.UserDetailResponse, error)

	// ModifyUserData now requires requester info for authorization logic.
	ModifyUserData(ctx context.Context, targetID int64, req dto.UpdateUserRequest, requesterID int64, requesterRole string) (*dto.UserDetailResponse, error)

	// DeactivateUserAccount now requires requester info to prevent self-deletion.
	DeactivateUserAccount(ctx context.Context, targetID int64, requesterID int64) error
}

type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// RegisterUser handles the registration of a new employee.
func (s *userService) RegisterUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserDetailResponse, error) {

	// 1. Check for Duplicate Data

	// Cheeck Username
	exists, err := s.userRepo.IsUsernameExists(ctx, req.Username, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		// [VIP FIX] Gunakan Sentinel Error
		return nil, response.ErrDuplicate
	}

	// Check Email
	exists, err = s.userRepo.IsEmailExists(ctx, req.Email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, response.ErrDuplicate
	}

	// Check Phone Number
	exists, err = s.userRepo.IsPhoneExists(ctx, req.PhoneNumber, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, response.ErrDuplicate
	}

	// 2. Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, response.ErrInternalServer // Hash gagal = Internal error
	}

	// 3. Prepare Model
	userModel := &models.User{
		FullName:     req.FullName,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		PhoneNumber:  req.PhoneNumber,
		IsActive:     true, // Default active upon creation
		CreatedAt:    time.Now(),
	}

	// 4. Insert into DB
	if err := s.userRepo.InsertUser(ctx, userModel); err != nil {
		return nil, err
	}

	// 5. Return Response
	return &dto.UserDetailResponse{
		ID:          userModel.ID,
		FullName:    userModel.FullName,
		Username:    userModel.Username,
		Email:       userModel.Email,
		Role:        userModel.Role,
		PhoneNumber: userModel.PhoneNumber,
		IsActive:    userModel.IsActive,
		CreatedAt:   userModel.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// RetrievedUserDirectory fetches a list of users with pagination and filters.
func (s *userService) GetUsers(ctx context.Context, page, perPage int, search, role string, status int) (*dto.UserListResponse, error) {

	// 1. Calculate Offset
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	// 2. Call Repository
	users, totalItems, err := s.userRepo.FetchUsers(ctx, perPage, offset, search, role, status)
	if err != nil {
		return nil, err
	}

	// 3. Map to DTO Summary
	// [OPTIMASI] Pre-allocate slice capacity
	userResponses := make([]dto.UserSummaryResponse, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, dto.UserSummaryResponse{
			ID:       u.ID,
			FullName: u.FullName,
			Username: u.Username,
			Role:     u.Role,
			IsActive: u.IsActive,
		})
	}

	// 4. Calculate Total Pages
	totalPages := 0
	if totalItems > 0 {
		totalPages = int((totalItems + int64(perPage) - 1) / int64(perPage))
	}

	return &dto.UserListResponse{
		Data: userResponses,
		Meta: response.MetaData{
			CurrentPage: page,
			PerPage:     perPage,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

// GetUserProfile retrieves detailed user information by ID.
func (s *userService) GetUserProfile(ctx context.Context, id int64) (*dto.UserDetailResponse, error) {

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Format Dates safely
	createdAtStr := user.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAtStr := ""
	if user.UpdatedAt != nil {
		updatedAtStr = user.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	lastLoginStr := ""
	if user.LastLoginAt != nil {
		lastLoginStr = user.LastLoginAt.Format("2006-01-02 15:04:05")
	}

	return &dto.UserDetailResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
		LastLoginAt: lastLoginStr,
		CreatedAt:   createdAtStr,
		UpdatedAt:   updatedAtStr,
	}, nil
}

// ModifyUserData updates user profile with strict security checks.
func (s *userService) ModifyUserData(ctx context.Context, targetID int64, req dto.UpdateUserRequest, requesterID int64, requesterRole string) (*dto.UserDetailResponse, error) {

	// 1. Retrieve Existing User
	existingUser, err := s.userRepo.FindByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	// 2. SECURITY GUARD: Access Control
	if requesterRole != "owner" {
		// Rule A: Non-owners can only edit their own profile
		if targetID != requesterID {
			// [VIP FIX] Gunakan Sentinel Error
			return nil, response.ErrForbidden
		}

		// Rule B: Non-owners CANNOT change Role or Active Status (Silent Ignore)
		req.Role = ""
		req.IsActive = nil
	}

	// 3. Update Fields (Partial Update Logic)

	if req.FullName != "" {
		existingUser.FullName = req.FullName
	}

	// [VIP FIX] Pengecekan Duplikat dengan Error Handling yang benar
	if req.Username != "" && req.Username != existingUser.Username {
		exists, err := s.userRepo.IsUsernameExists(ctx, req.Username, targetID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, response.ErrDuplicate
		}
		existingUser.Username = req.Username
	}

	if req.Email != "" && req.Email != existingUser.Email {
		exists, err := s.userRepo.IsEmailExists(ctx, req.Email, targetID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, response.ErrDuplicate
		}
		existingUser.Email = req.Email
	}

	if req.PhoneNumber != "" && req.PhoneNumber != existingUser.PhoneNumber {
		exists, err := s.userRepo.IsPhoneExists(ctx, req.PhoneNumber, targetID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, response.ErrDuplicate
		}
		existingUser.PhoneNumber = req.PhoneNumber
	}

	// Role is only updated if it passed the Security Guard above
	if req.Role != "" {
		existingUser.Role = req.Role
	}

	if req.Password != "" {
		hashedPwd, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, response.ErrInternalServer
		}
		existingUser.PasswordHash = hashedPwd
	}

	if req.IsActive != nil {
		existingUser.IsActive = *req.IsActive
	}

	// 4. Update Timestamp
	now := time.Now()
	existingUser.UpdatedAt = &now

	// 5. Save Changes
	if err := s.userRepo.UpdateUser(ctx, existingUser); err != nil {
		return nil, err
	}

	// 6. Return Updated Profile
	return s.GetUserProfile(ctx, targetID)
}

// DeactivateUserAccount handles soft deletion of a user.
func (s *userService) DeactivateUserAccount(ctx context.Context, targetID int64, requesterID int64) error {

	// 1. SECURITY GUARD: Anti Self-Deletion
	if targetID == requesterID {
		// [VIP FIX] Gunakan Sentinel Error
		return response.ErrForbidden
	}

	// 2. Check if user exists
	_, err := s.userRepo.FindByID(ctx, targetID)
	if err != nil {
		return err
	}

	// 3. Execute Soft Delete
	return s.userRepo.DeleteUser(ctx, targetID)
}
