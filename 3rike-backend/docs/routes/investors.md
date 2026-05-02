# /api/investors — Investor Profiles

All routes require `Authorization: Bearer <token>`.

An investor buys fractions of tokenized tricycles and earns weekly yield.

---

## Create investor profile

```http
POST /api/investors
Authorization: Bearer <token>
```

```json
{
  "full_name": "Amara Nwosu",
  "wallet_address": "0xabc123..."
}
```

**Response 201**
```json
{
  "id": 1,
  "user_id": 5,
  "full_name": "Amara Nwosu",
  "wallet_address": "0xabc123...",
  "created_at": "2026-05-01T10:00:00Z"
}
```

---

## Get investor by ID

```http
GET /api/investors/:id
Authorization: Bearer <token>
```

---

## List all investors

```http
GET /api/investors
Authorization: Bearer <token>
```

---

## List investor's investments (fractions)

```http
GET /api/investors/:id/investments
Authorization: Bearer <token>
```

Returns all tricycle fractions owned by this investor.

**Response 200**
```json
[
  {
    "id": 1,
    "tricycle_id": 2,
    "investor_id": 1,
    "units": 10,
    "price_per_unit": 18.00,
    "created_at": "2026-05-01T10:00:00Z"
  }
]
```
