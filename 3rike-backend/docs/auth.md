# Authentication & Sessions

## Overview

3riKE uses JWT tokens backed by Redis session storage. Each user can have **max 3 concurrent sessions** — the oldest is evicted automatically on the 4th login.

## Register

```http
POST /auth/register
Content-Type: application/json

{
  "email": "driver@example.com",
  "password": "securepassword",
  "role": "driver"          // "driver" | "investor" | "admin"
}
```

Response:
```json
{
  "id": 1,
  "email": "driver@example.com",
  "role": "driver",
  "created_at": "2026-05-01T10:00:00Z"
}
```

## Login

```http
POST /auth/login
Content-Type: application/json

{
  "email": "driver@example.com",
  "password": "securepassword"
}
```

Response:
```json
{
  "token": "eyJhbGci...",
  "session_id": "1-1777551649000000000"
}
```

Use the token in all subsequent requests:
```http
Authorization: Bearer eyJhbGci...
```

## Logout

```http
POST /auth/logout
Authorization: Bearer <token>
```

Immediately invalidates the session in Redis — the token is dead even if not expired.

## Session management

```http
# List all active sessions (devices)
GET /auth/sessions
Authorization: Bearer <token>

# Revoke a specific session (log out another device)
DELETE /auth/sessions/:sessionID
Authorization: Bearer <token>
```

## User profile

```http
# Get current user
GET /auth/me

# Update email
PUT /auth/profile
{ "email": "new@email.com" }

# Change password
PUT /auth/password
{ "old_password": "...", "new_password": "..." }

# Delete account
DELETE /auth/account
```

## Link Canton wallet

After getting your Canton party ID from the wallet UI:

```http
PUT /auth/wallet
Authorization: Bearer <token>

{
  "canton_party_id": "abc123::1220xyz..."
}
```

Once linked, all your contract actions (tokenize, fractionalize) are signed by your own Canton party.

## Get CC balance

```http
GET /auth/wallet/balance
Authorization: Bearer <token>
```

Returns your Canton Coin balance live from the ledger:
```json
{
  "round": 42600,
  "effective_unlocked_qty": "646.2453147215",
  "effective_locked_qty": "0.0000000000",
  "total_holding_fees": "0.0360255220"
}
```

## JWT claims

The JWT contains:
```json
{
  "sub": 1,                          // user DB ID
  "role": "driver",
  "session_id": "1-177755...",
  "canton_party_id": "abc::1220...", // set after wallet link
  "exp": 1777640000
}
```
