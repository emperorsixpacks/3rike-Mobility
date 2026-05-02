# /api/tricycles — Tricycle Assets

All routes require `Authorization: Bearer <token>`.

Tricycles go through this lifecycle:
```
available → tokenized → fractionalized
```

---

## Register a tricycle

```http
POST /api/tricycles
Authorization: Bearer <token>
```

```json
{
  "driver_id": 1,
  "make": "Bajaj",
  "model": "RE 4S",
  "is_ev": false,
  "price_usd": 1800.00
}
```

**Response 201**
```json
{
  "id": 1,
  "driver_id": 1,
  "make": "Bajaj",
  "model": "RE 4S",
  "is_ev": false,
  "price_usd": 1800.00,
  "status": "available",
  "contract_id": "",
  "total_fractions": 0,
  "created_at": "2026-05-01T10:00:00Z"
}
```

---

## Get tricycle by ID

```http
GET /api/tricycles/:id
Authorization: Bearer <token>
```

---

## List all tricycles

```http
GET /api/tricycles
Authorization: Bearer <token>
```

---

## Tokenize a tricycle

Creates a `TricycleToken` contract on the Canton ledger. The caller's Canton party ID (from JWT) is used as the signatory.

```http
POST /api/tricycles/:id/tokenize
Authorization: Bearer <token>
```

No request body needed.

**Response 200**
```json
{
  "id": 1,
  "status": "tokenized",
  "contract_id": "00ae971d4a66961fed206705169c7a328a4456f4...",
  ...
}
```

> Requires `canton_party_id` to be set on the user (via `PUT /auth/wallet`).
> Falls back to operator party if not set.

---

## Fractionalize a tricycle

Exercises the `Fractionalize` choice on the Canton contract, splitting the tricycle into N investor shares.

```http
POST /api/tricycles/:id/fractionalize
Authorization: Bearer <token>
```

```json
{
  "total_fractions": 100
}
```

**Response 200**
```json
{
  "id": 1,
  "status": "fractionalized",
  "total_fractions": 100,
  "contract_id": "0051fccf079c01fe968b04d18841db04c2d96c61...",
  ...
}
```

> Tricycle must be `tokenized` before fractionalization.
