# API Documentation

Base URL: `http://localhost:8080`

## Response Format

All endpoints return JSON responses in the following format:

**Success Response:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {}
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error information"
}
```

---

## Health Check

### Check API Health
```
GET /health
```

**Response:**
```
OK
```

---

## Borrowers

### Create Borrower
```
POST /api/v1/borrowers
```

**Request Body:**
```json
{
  "borrower_id_number": "ID123456",
  "full_name": "John Doe",
  "email": "john@example.com",
  "phone": "+628123456789",
  "address": "123 Main St, Jakarta"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Borrower created successfully",
  "data": {
    "id": 1,
    "borrower_id_number": "ID123456",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+628123456789",
    "address": "123 Main St, Jakarta",
    "created_at": "2025-11-19T00:00:00Z",
    "updated_at": "2025-11-19T00:00:00Z"
  }
}
```

### Get Borrower by ID
```
GET /api/v1/borrowers/{id}
```

**Path Parameters:**
- `id` (integer, required): Borrower ID

**Response:**
```json
{
  "success": true,
  "message": "Borrower retrieved successfully",
  "data": {
    "id": 1,
    "borrower_id_number": "ID123456",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+628123456789",
    "address": "123 Main St, Jakarta",
    "created_at": "2025-11-19T00:00:00Z",
    "updated_at": "2025-11-19T00:00:00Z"
  }
}
```

### Update Borrower
```
PUT /api/v1/borrowers/{id}
```

**Path Parameters:**
- `id` (integer, required): Borrower ID

**Request Body:**
```json
{
  "borrower_id_number": "ID123456",
  "full_name": "John Doe Updated",
  "email": "john.updated@example.com",
  "phone": "+628123456789",
  "address": "456 New St, Jakarta"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Borrower updated successfully",
  "data": {
    "id": 1,
    "borrower_id_number": "ID123456",
    "full_name": "John Doe Updated",
    "email": "john.updated@example.com",
    "phone": "+628123456789",
    "address": "456 New St, Jakarta",
    "created_at": "2025-11-19T00:00:00Z",
    "updated_at": "2025-11-19T00:00:00Z"
  }
}
```

### Delete Borrower
```
DELETE /api/v1/borrowers/{id}
```

**Path Parameters:**
- `id` (integer, required): Borrower ID

**Response:**
```json
{
  "success": true,
  "message": "Borrower deleted successfully",
  "data": null
}
```

### List Borrowers
```
GET /api/v1/borrowers?offset=0&limit=10
```

**Query Parameters:**
- `offset` (integer, optional, default: 0): Number of records to skip
- `limit` (integer, optional, default: 10): Maximum number of records to return

**Response:**
```json
{
  "success": true,
  "message": "Borrowers retrieved successfully",
  "data": [
    {
      "id": 1,
      "borrower_id_number": "ID123456",
      "full_name": "John Doe",
      "email": "john@example.com",
      "phone": "+628123456789",
      "address": "123 Main St, Jakarta",
      "created_at": "2025-11-19T00:00:00Z",
      "updated_at": "2025-11-19T00:00:00Z"
    }
  ]
}
```

---

## Loans

### Create Loan
```
POST /api/v1/loans
```

**Request Body:**
```json
{
  "borrower_id": 1,
  "principal_amount": 10000000,
  "rate": 12.5,
  "roi": 15.0,
  "agreement_letter_link": "https://storage.example.com/agreement.pdf"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Loan created successfully",
  "data": {
    "id": 1,
    "loan_id": "LOAN-20251119-001",
    "borrower_id": 1,
    "principal_amount": 10000000,
    "rate": 12.5,
    "roi": 15.0,
    "agreement_letter_link": "https://storage.example.com/agreement.pdf",
    "current_state": "proposed",
    "total_invested_amount": 0,
    "created_at": "2025-11-19T00:00:00Z",
    "updated_at": "2025-11-19T00:00:00Z"
  }
}
```

### Get Loan by ID
```
GET /api/v1/loans/{id}
```

**Path Parameters:**
- `id` (integer, required): Loan ID

**Response:**
```json
{
  "success": true,
  "message": "Loan retrieved successfully",
  "data": {
    "id": 1,
    "loan_id": "LOAN-20251119-001",
    "borrower_id": 1,
    "principal_amount": 10000000,
    "rate": 12.5,
    "roi": 15.0,
    "agreement_letter_link": "https://storage.example.com/agreement.pdf",
    "current_state": "proposed",
    "total_invested_amount": 0,
    "created_at": "2025-11-19T00:00:00Z",
    "updated_at": "2025-11-19T00:00:00Z"
  }
}
```

### Update Loan
```
PUT /api/v1/loans/{id}
```

**Path Parameters:**
- `id` (integer, required): Loan ID

**Request Body:**
```json
{
  "borrower_id": 1,
  "principal_amount": 12000000,
  "rate": 13.0,
  "roi": 16.0,
  "agreement_letter_link": "https://storage.example.com/agreement-updated.pdf"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Loan updated successfully",
  "data": {
    "borrower_id": 1,
    "principal_amount": 12000000,
    "rate": 13.0,
    "roi": 16.0,
    "agreement_letter_link": "https://storage.example.com/agreement-updated.pdf"
  }
}
```

### Delete Loan
```
DELETE /api/v1/loans/{id}
```

**Path Parameters:**
- `id` (integer, required): Loan ID

**Response:**
```json
{
  "success": true,
  "message": "Loan deleted successfully",
  "data": null
}
```

### List Loans
```
GET /api/v1/loans?state=proposed&offset=0&limit=10
```

**Query Parameters:**
- `state` (string, optional): Filter by loan state (proposed, approved, invested, disbursed)
- `offset` (integer, optional, default: 0): Number of records to skip
- `limit` (integer, optional, default: 10): Maximum number of records to return

**Response:**
```json
{
  "success": true,
  "message": "Loans retrieved successfully",
  "data": [
    {
      "id": 1,
      "loan_id": "LOAN-20251119-001",
      "borrower_id": 1,
      "principal_amount": 10000000,
      "rate": 12.5,
      "roi": 15.0,
      "agreement_letter_link": "https://storage.example.com/agreement.pdf",
      "current_state": "proposed",
      "total_invested_amount": 0,
      "created_at": "2025-11-19T00:00:00Z",
      "updated_at": "2025-11-19T00:00:00Z"
    }
  ]
}
```

### Approve Loan
```
POST /api/v1/loans/{id}/approve
```

**Path Parameters:**
- `id` (integer, required): Loan ID

**Request Body:**
```json
{
  "field_validator_employee_id": "EMP001",
  "proof_image_url": "https://storage.example.com/proof.jpg"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Loan approved successfully",
  "data": null
}
```

**State Transition:** `proposed` → `approved`

### Invest in Loan
```
POST /api/v1/loans/{id}/invest
```

**Path Parameters:**
- `id` (integer, required): Loan ID

**Request Body:**
```json
{
  "investor_id": 1,
  "investment_amount": 5000000
}
```

**Response:**
```json
{
  "success": true,
  "message": "Investment completed successfully",
  "data": null
}
```

**State Transition:** `approved` → `invested` (when total investments >= principal amount)

**Notes:**
- Multiple investors can invest in the same loan
- Loan state changes to "invested" when total invested amount reaches or exceeds principal amount
- Investors receive email notifications upon successful investment

### Disburse Loan
```
POST /api/v1/loans/{id}/disburse
```

**Path Parameters:**
- `id` (integer, required): Loan ID

**Request Body:**
```json
{
  "field_officer_employee_id": "EMP002",
  "agreement_letter_signed_url": "https://storage.example.com/agreement-signed.pdf"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Loan disbursed successfully",
  "data": null
}
```

**State Transition:** `invested` → `disbursed`

### Get Loans by State
```
GET /api/v1/loans/state/{state}
```

**Path Parameters:**
- `state` (string, required): Loan state (proposed, approved, invested, disbursed)

**Response:**
```json
{
  "success": true,
  "message": "Loans retrieved successfully",
  "data": [
    {
      "id": 1,
      "loan_id": "LOAN-20251119-001",
      "borrower_id": 1,
      "principal_amount": 10000000,
      "rate": 12.5,
      "roi": 15.0,
      "agreement_letter_link": "https://storage.example.com/agreement.pdf",
      "current_state": "proposed",
      "total_invested_amount": 0,
      "created_at": "2025-11-19T00:00:00Z",
      "updated_at": "2025-11-19T00:00:00Z"
    }
  ]
}
```

---

## Investors

### Create Investor
```
POST /api/v1/investors
```

**Request Body:**
```json
{
  "investor_id": "INV001",
  "full_name": "Jane Smith",
  "email": "jane@example.com",
  "phone": "+628987654321"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Investor created successfully",
  "data": {
    "id": 1,
    "investor_id": "INV001",
    "full_name": "Jane Smith",
    "email": "jane@example.com",
    "phone": "+628987654321",
    "created_at": "2025-11-19T00:00:00Z",
    "updated_at": "2025-11-19T00:00:00Z"
  }
}
```

### Get Investor by ID
```
GET /api/v1/investors/{id}
```

**Path Parameters:**
- `id` (integer, required): Investor ID

**Response:**
```json
{
  "success": true,
  "message": "Investor retrieved successfully",
  "data": {
    "id": 1,
    "investor_id": "INV001",
    "full_name": "Jane Smith",
    "email": "jane@example.com",
    "phone": "+628987654321",
    "created_at": "2025-11-19T00:00:00Z",
    "updated_at": "2025-11-19T00:00:00Z"
  }
}
```

### Update Investor
```
PUT /api/v1/investors/{id}
```

**Path Parameters:**
- `id` (integer, required): Investor ID

**Request Body:**
```json
{
  "investor_id": "INV001",
  "full_name": "Jane Smith Updated",
  "email": "jane.updated@example.com",
  "phone": "+628987654321"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Investor updated successfully",
  "data": {
    "id": 1,
    "investor_id": "INV001",
    "full_name": "Jane Smith Updated",
    "email": "jane.updated@example.com",
    "phone": "+628987654321",
    "created_at": "2025-11-19T00:00:00Z",
    "updated_at": "2025-11-19T00:00:00Z"
  }
}
```

### Delete Investor
```
DELETE /api/v1/investors/{id}
```

**Path Parameters:**
- `id` (integer, required): Investor ID

**Response:**
```json
{
  "success": true,
  "message": "Investor deleted successfully",
  "data": null
}
```

### List Investors
```
GET /api/v1/investors?offset=0&limit=10
```

**Query Parameters:**
- `offset` (integer, optional, default: 0): Number of records to skip
- `limit` (integer, optional, default: 10): Maximum number of records to return

**Response:**
```json
{
  "success": true,
  "message": "Investors retrieved successfully",
  "data": [
    {
      "id": 1,
      "investor_id": "INV001",
      "full_name": "Jane Smith",
      "email": "jane@example.com",
      "phone": "+628987654321",
      "created_at": "2025-11-19T00:00:00Z",
      "updated_at": "2025-11-19T00:00:00Z"
    }
  ]
}
```

---

## Loan State Transitions

The loan lifecycle follows a strict state machine:

```
proposed → approved → invested → disbursed
```

**State Descriptions:**

1. **proposed**: Initial state when a loan is created
2. **approved**: Loan has been validated and approved by a field validator
3. **invested**: Loan has received sufficient investment (total invested >= principal amount)
4. **disbursed**: Loan funds have been disbursed to the borrower

**Rules:**
- State transitions can only move forward, never backward
- Each transition requires specific data and validation
- State history is tracked in the `loan_state_history` table

---

## Error Codes

| HTTP Status | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid input data |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error |

---

## Example Usage

### Complete Loan Workflow

1. **Create a borrower:**
```bash
curl -X POST http://localhost:8080/api/v1/borrowers \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id_number": "ID123456",
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+628123456789",
    "address": "123 Main St, Jakarta"
  }'
```

2. **Create a loan:**
```bash
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": 1,
    "principal_amount": 10000000,
    "rate": 12.5,
    "roi": 15.0,
    "agreement_letter_link": "https://storage.example.com/agreement.pdf"
  }'
```

3. **Approve the loan:**
```bash
curl -X POST http://localhost:8080/api/v1/loans/1/approve \
  -H "Content-Type: application/json" \
  -d '{
    "field_validator_employee_id": "EMP001",
    "proof_image_url": "https://storage.example.com/proof.jpg"
  }'
```

4. **Create an investor:**
```bash
curl -X POST http://localhost:8080/api/v1/investors \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": "INV001",
    "full_name": "Jane Smith",
    "email": "jane@example.com",
    "phone": "+628987654321"
  }'
```

5. **Invest in the loan:**
```bash
curl -X POST http://localhost:8080/api/v1/loans/1/invest \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": 1,
    "investment_amount": 10000000
  }'
```

6. **Disburse the loan:**
```bash
curl -X POST http://localhost:8080/api/v1/loans/1/disburse \
  -H "Content-Type: application/json" \
  -d '{
    "field_officer_employee_id": "EMP002",
    "agreement_letter_signed_url": "https://storage.example.com/agreement-signed.pdf"
  }'
```
