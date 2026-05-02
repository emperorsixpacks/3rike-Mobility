# /api/yield — Investor Yield Payouts

All routes require `Authorization: Bearer <token>`.

Investors earn weekly yield from the tricycles they've invested in. Yield is distributed by the platform operator each week based on driver repayments.

---

## Get yield payouts for an investor

```http
GET /api/yield/investor/:investorID
Authorization: Bearer <token>
```

**Response 200**
```json
[
  {
    "id": 1,
    "investor_id": 1,
    "fraction_id": 3,
    "amount_usdc": 2.50,
    "week_number": 12,
    "created_at": "2026-05-01T10:00:00Z"
  },
  {
    "id": 2,
    "investor_id": 1,
    "fraction_id": 3,
    "amount_usdc": 2.50,
    "week_number": 13,
    "created_at": "2026-05-08T10:00:00Z"
  }
]
```

Each payout is linked to a specific `fraction_id` — the investor's share in a tricycle.
