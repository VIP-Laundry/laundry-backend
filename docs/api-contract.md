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

- POST /api/auth/login

- POST /api/auth/logout

- GET /api/auth/me

### Users

- POST /api/users

- GET /api/users

- GET /api/users/{id}

- PUT /api/users/{id}

- DELETE /api/users/{id}

### Service Categories

- POST /api/categories

- GET /api/categories

- GET /api/categories/{id}

- PUT /api/categories/{id}

- DELETE /api/categories/{id}

### Services

- POST /api/services

- GET /api/services

- GET /api/services/{id}

- PUT /api/services/{id}

- DELETE /api/services/{id}

### Orders

- POST /api/orders

- GET /api/orders

- GET /api/orders/{id}

- PUT /api/orders/{id}

- PATCH /api/orders/{id}

### Payments

- GET /api/payments

- GET /api/payments/{id}

- PATCH /api/payments/{id}

### Deliveries

- GET /api/deliveries

- GET /api/deliveries/{id}

- GET /api/deliveries/my-tasks

- PATCH /api/deliveries/{id}

### Customer (Endpoint Public Tracking)

- GET /api/orders/track/{invoice_number}

### Reports

- GET /api/reports/dashboard

- GET /api/reports/revenue

- GET /api/reports/payments

- GET /api/reports/employees
