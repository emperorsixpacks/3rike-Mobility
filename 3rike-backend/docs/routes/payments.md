# /api/payments — Driver Weekly Repayments

All routes require `Authorization: Bearer <token>`.

Drivers make weekly payments in local currency. Each confirmed payment decrements `weeks_remaining` on the driver profile by 1 (out of 70).

---

## Record a payment

```http
POST /api/payments
Authorization: Bearer <token>
```

```json
{
  "driver_id": 1,
  "tricycle_id": 1,
  "amount_local": 15000.00,
  "amount_usdc": 10.50,
  "currency": "NGN",
  "week_number": 12
}
```

**Response 201**
```json
{
  "id": 1,
  "driver_id": 1,
  "tricycle_id": 1,
  "amount_local": 15000.00,
  "amount_usdc": 10.50,
  "currency": "NGN",
  "status": "confirmed",
  "week_number": 12,
  "created_at": "2026-05-01T10:00:00Z"
}
```

---

## Get payments for a driver

```http
GET /api/payments/driver/:driverID
Authorization: Bearer <token>
```

**Response 200** — array of payment objects.
