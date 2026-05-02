# /api/savings — Driver Savings Account

All routes require `Authorization: Bearer <token>`.

Drivers can save on the platform and earn from a shared interest pool. One savings account per driver, auto-created on first deposit.

---

## Deposit

```http
POST /api/savings/deposit
Authorization: Bearer <token>
```

```json
{
  "driver_id": 1,
  "amount_usdc": 25.00
}
```

**Response 200**
```json
{
  "id": 1,
  "driver_id": 1,
  "balance_usdc": 75.00,
  "created_at": "2026-05-01T10:00:00Z"
}
```

---

## Get balance

```http
GET /api/savings/:driverID/balance
Authorization: Bearer <token>
```

**Response 200**
```json
{
  "id": 1,
  "driver_id": 1,
  "balance_usdc": 75.00,
  "created_at": "2026-05-01T10:00:00Z"
}
```
