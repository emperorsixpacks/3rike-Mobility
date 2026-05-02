# /api/loans — Driver Credit Facility

All routes require `Authorization: Bearer <token>`.

Drivers with a credit score ≥ 500 and good repayment history can access loans. Repayments are structured weekly.

---

## Apply for a loan

```http
POST /api/loans
Authorization: Bearer <token>
```

```json
{
  "driver_id": 1,
  "principal_usdc": 200.00,
  "weekly_repayment": 10.00
}
```

Returns `422` if credit score < 500.

**Response 201**
```json
{
  "id": 1,
  "driver_id": 1,
  "principal_usdc": 200.00,
  "remaining_usdc": 200.00,
  "weekly_repayment": 10.00,
  "status": "active",
  "created_at": "2026-05-01T10:00:00Z"
}
```

---

## Get loan by ID

```http
GET /api/loans/:id
Authorization: Bearer <token>
```

---

## Make a repayment

```http
PUT /api/loans/:id/repay
Authorization: Bearer <token>
```

```json
{ "amount_usdc": 10.00 }
```

**Response 200**
```json
{
  "id": 1,
  "remaining_usdc": 190.00,
  "status": "active",
  ...
}
```

When `remaining_usdc` reaches 0, `status` becomes `"repaid"`.
