# Loan Engine API Endpoints

## Overview
This document outlines the RESTful API endpoints for the loan engine system, following standard REST conventions and supporting the loan lifecycle management.

## Authentication
All endpoints (except user registration and login) require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

## API Versioning
All endpoints are versioned using path prefix: `/api/v1/`

## Common Response Format

### Success Response
```json
{
  "success": true,
  "data": {},
  "message": "Success message"
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  }
}
```

## Endpoints

### 1. Borrower Management

#### Create Borrower
- **POST** `/api/v1/borrowers`
- **Description**: Create a new borrower
- **Request Body**:
```json
{
  "borrower_id_number": "string",
  "full_name": "string",
  "email": "string",
  "phone": "string",
  "address": "string"
}
```
- **Response**: Created borrower object
- **Authentication**: Admin/Staff

#### Get Borrower
- **GET** `/api/v1/borrowers/{id}`
- **Description**: Get borrower details by ID
- **Response**: Borrower object
- **Authentication**: Admin/Staff

#### Update Borrower
- **PUT** `/api/v1/borrowers/{id}`
- **Description**: Update borrower details
- **Request Body**: Same as create borrower
- **Response**: Updated borrower object
- **Authentication**: Admin/Staff

### 2. Loan Management

#### Create Loan
- **POST** `/api/v1/loans`
- **Description**: Create a new loan (initial state: proposed)
- **Request Body**:
```json
{
  "borrower_id": "integer",
  "principal_amount": "decimal",
  "rate": "decimal",
  "roi": "decimal",
  "agreement_letter_link": "string"
}
```
- **Response**: Created loan object
- **Authentication**: Admin/Staff

#### Get Loan
- **GET** `/api/v1/loans/{id}`
- **Description**: Get loan details by ID
- **Response**: Loan object with all related information
- **Authentication**: Admin/Staff/Investor (limited access)

#### Get Loans by State
- **GET** `/api/v1/loans?state={state}&page={page}&limit={limit}`
- **Description**: Get loans filtered by state with pagination
- **Response**: Paginated list of loans
- **Authentication**: Admin/Staff

#### Get All Loans
- **GET** `/api/v1/loans?page={page}&limit={limit}`
- **Description**: Get all loans with pagination
- **Response**: Paginated list of loans
- **Authentication**: Admin/Staff

### 3. Loan Approval

#### Approve Loan
- **POST** `/api/v1/loans/{id}/approve`
- **Description**: Approve a loan and transition to approved state
- **Request Body**:
```json
{
  "field_validator_employee_id": "string",
  "proof_image_url": "string"
}
```
- **Response**: Updated loan object
- **Authentication**: Staff
- **Business Rules**:
  - Loan must be in 'proposed' state
  - All required approval information must be provided

#### Get Loan Approval
- **GET** `/api/v1/loans/{id}/approval`
- **Description**: Get approval details for a loan
- **Response**: Loan approval object
- **Authentication**: Admin/Staff

### 4. Loan Investment

#### Create Investment
- **POST** `/api/v1/loans/{id}/invest`
- **Description**: Create an investment in a loan
- **Request Body**:
```json
{
  "investor_id": "integer",
  "investment_amount": "decimal"
}
```
- **Response**: Investment object
- **Authentication**: Investor
- **Business Rules**:
  - Loan must be in 'approved' state
  - Investment amount must not exceed remaining principal
  - Total investment cannot exceed principal amount

#### Get Loan Investments
- **GET** `/api/v1/loans/{id}/investments`
- **Description**: Get all investments for a loan
- **Response**: List of investment objects
- **Authentication**: Admin/Staff/Investor (limited access)

#### Get Investor Investments
- **GET** `/api/v1/investors/{id}/investments`
- **Description**: Get all investments by an investor
- **Response**: List of investment objects
- **Authentication**: Admin/Staff/Investor (own investments only)

### 5. Loan Disbursement

#### Disburse Loan
- **POST** `/api/v1/loans/{id}/disburse`
- **Description**: Disburse a loan and transition to disbursed state
- **Request Body**:
```json
{
  "field_officer_employee_id": "string",
  "agreement_letter_signed_url": "string"
}
```
- **Response**: Updated loan object
- **Authentication**: Staff
- **Business Rules**:
  - Loan must be in 'invested' state
  - Total investment must equal principal amount
  - All required disbursement information must be provided

#### Get Loan Disbursement
- **GET** `/api/v1/loans/{id}/disbursement`
- **Description**: Get disbursement details for a loan
- **Response**: Loan disbursement object
- **Authentication**: Admin/Staff

### 6. Investor Management

#### Create Investor
- **POST** `/api/v1/investors`
- **Description**: Create a new investor
- **Request Body**:
```json
{
  "investor_id": "string",
  "full_name": "string",
  "email": "string",
  "phone": "string"
}
```
- **Response**: Created investor object
- **Authentication**: Admin

#### Get Investor
- **GET** `/api/v1/investors/{id}`
- **Description**: Get investor details by ID
- **Response**: Investor object
- **Authentication**: Admin/Staff/Investor (own details only)

#### Update Investor
- **PUT** `/api/v1/investors/{id}`
- **Description**: Update investor details
- **Request Body**: Same as create investor
- **Response**: Updated investor object
- **Authentication**: Admin/Staff/Investor (own details only)

### 7. User Authentication

#### Register User
- **POST** `/api/v1/auth/register`
- **Description**: Register a new user
- **Request Body**:
```json
{
  "email": "string",
  "password": "string",
  "user_type": "string",
  "full_name": "string"
}
```
- **Response**: User object with JWT token
- **Authentication**: None

#### Login User
- **POST** `/api/v1/auth/login`
- **Description**: Authenticate user and return JWT token
- **Request Body**:
```json
{
  "email": "string",
  "password": "string"
}
```
- **Response**: JWT token
- **Authentication**: None

#### Refresh Token
- **POST** `/api/v1/auth/refresh`
- **Description**: Refresh JWT token
- **Request Body**:
```json
{
  "refresh_token": "string"
}
```
- **Response**: New JWT token
- **Authentication**: None

### 8. Loan State History

#### Get Loan State History
- **GET** `/api/v1/loans/{id}/state-history`
- **Description**: Get all state transitions for a loan
- **Response**: List of state transition objects
- **Authentication**: Admin/Staff

## Error Codes

- `LOAN_NOT_FOUND`: Loan with specified ID does not exist
- `BORROWER_NOT_FOUND`: Borrower with specified ID does not exist
- `INVESTOR_NOT_FOUND`: Investor with specified ID does not exist
- `INVALID_STATE_TRANSITION`: Attempted to perform invalid state transition
- `INSUFFICIENT_FUNDS`: Investment amount exceeds available principal
- `LOAN_ALREADY_APPROVED`: Loan is already in approved state
- `LOAN_ALREADY_INVESTED`: Loan is already in invested state
- `LOAN_ALREADY_DISBURSED`: Loan is already in disbursed state
- `INVALID_APPROVAL_DATA`: Missing required approval information
- `INVALID_DISBURSEMENT_DATA`: Missing required disbursement information
- `DUPLICATE_INVESTMENT`: Investor already invested in this loan
- `UNAUTHORIZED_ACCESS`: User does not have permission to access resource
- `VALIDATION_ERROR`: Request data validation failed
- `INTERNAL_ERROR`: Internal server error occurred