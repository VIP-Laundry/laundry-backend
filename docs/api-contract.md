# API Contract - Laundry Management System

## General Rules

- Base URL: /api/v1
- All requests and responses use JSON
- Authentication via Authorization header
- Timestamps use ISO 8601 format

## Authentication

Authenticated endpoints require:
Authorization: Bearer <token>

## Roles:

- owner
- kasir
- staff
- courier

## Roles & Access Rules

- Owner: read-only access to summaries
- Kasir: create orders, manage payments
- Staff: update production status
- Courier: handle delivery status
- Customer: public access (no auth) for order status lookup

## Endpoints

### Auths

- POST /api/v1/auth/login

- POST /api/v1/auth/refresh-token

- POST /api/v1/auth/logout

- GET /api/v1/auth/me

### Users

- POST /api/v1/users

- GET /api/v1/users

- GET /api/v1/users/{id}

- PUT /api/v1/users/{id}

- DELETE /api/v1/users/{id}

### Service Categories

- POST /api/v1/categories

- GET /api/v1/categories

- GET /api/v1/categories/{id}

- PUT /api/v1/categories/{id}

- DELETE /api/v1/categories/{id}

### Services

- POST /api/v1/services

- GET /api/v1/services

- GET /api/v1/services/{id}

- PUT /api/v1/services/{id}

- DELETE /api/v1/services/{id}

### Orders

- POST /api/v1/orders

- GET /api/v1/orders

- GET /api/v1/orders/{id}

- PUT /api/v1/orders/{id}

- PATCH /api/v1/orders/{id}

### Payments

- GET /api/v1/payments

- GET /api/v1/payments/{id}

- PATCH /api/v1/payments/{id}

### Deliveries

- GET /api/v1/deliveries

- GET /api/v1/deliveries/{id}

- GET /api/v1/deliveries/my-tasks

- PATCH /api/v1/deliveries/{id}

### Customer (Endpoint Public Tracking)

- GET /api/v1/orders/track/{invoice_number}

### Reports

- GET /api/v1/reports/dashboard

- GET /api/v1/reports/revenue

- GET /api/v1/reports/payments

- GET /api/v1/reports/employees
