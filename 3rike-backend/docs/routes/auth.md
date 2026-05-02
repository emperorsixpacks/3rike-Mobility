# /auth — Authentication & User Management

Base: `POST /auth/register`, `POST /auth/login` are **public**.
Everything else requires `Authorization: Bearer <token>`.

---

## Register

```http
POST /auth/register
```

```json
{
  "email": "driver@example.com",
  "password": "secret123",
  "role": "driver"
}
```
`role` is one of: `driver` | `investor` | `admin`

**Response 201**
```json
{
  "id": 1,
  "email": "driver@example.com",
  "role": "driver",
  "created_at": "2026-05-01T10:00:00Z"
}
```

---

## Login

```http
POST /auth/login
```

```json
{
  "email": "driver@example.com",
  "password": "secret123"
}
```

**Response 200**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "session_id": "1-1777551649000000000"
}
```

Use `token` in all subsequent requests as `Authorization: Bearer <token>`.

---

## Logout

```http
POST /auth/logout
Authorization: Bearer <token>
```

Kills the current session immediately — token is invalid even before expiry.

**Response 200**
```json
{ "message": "logged out" }
```

---

## Get current user

```http
GET /auth/me
Authorization: Bearer <token>
```

**Response 200**
```json
{
  "id": 1,
  "email": "driver@example.com",
  "role": "driver",
  "canton_party_id": "abc::1220...",
  "created_at": "2026-05-01T10:00:00Z"
}
```

---

## Update email / profile

```http
PUT /auth/profile
Authorization: Bearer <token>
```

```json
{ "email": "newemail@example.com" }
```

---

## Change password

```http
PUT /auth/password
Authorization: Bearer <token>
```

```json
{
  "old_password": "secret123",
  "new_password": "newsecret456"
}
```

---

## List active sessions

```http
GET /auth/sessions
Authorization: Bearer <token>
```

Returns all active sessions (devices) for the current user. Max 3.

**Response 200**
```json
[
  {
    "id": "1-1777551649000000000",
    "user_id": 1,
    "role": "driver",
    "created_at": "2026-05-01T10:00:00Z",
    "expires_at": "2026-05-04T10:00:00Z"
  }
]
```

---

## Revoke a session (log out another device)

```http
DELETE /auth/sessions/:sessionID
Authorization: Bearer <token>
```

---

## Link Canton wallet

After getting your party ID from the Canton wallet UI:

```http
PUT /auth/wallet
Authorization: Bearer <token>
```

```json
{
  "canton_party_id": "70af8ee8-8bc6-4a66-81a6-d375a678273a::1220195a..."
}
```

**Response 200** — updated user object with `canton_party_id` set.

---

## Get CC (Canton Coin) balance

```http
GET /auth/wallet/balance
Authorization: Bearer <token>
```

Fetches live from the Canton ledger using your token.

**Response 200**
```json
{
  "round": 42600,
  "effective_unlocked_qty": "646.2453147215",
  "effective_locked_qty": "0.0000000000",
  "total_holding_fees": "0.0360255220"
}
```

---

## Delete account

```http
DELETE /auth/account
Authorization: Bearer <token>
```

**Response 204** — no content.
