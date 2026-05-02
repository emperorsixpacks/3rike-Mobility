# /waitlist — Pre-launch Waitlist

All routes are **public** — no auth required.

---

## Join the waitlist

```http
POST /waitlist/join
```

```json
{
  "email": "user@example.com",
  "phone": "+2348012345678",
  "referral_code": "ABC123"
}
```

`phone` and `referral_code` are optional. If `referral_code` matches an existing entry, the referrer moves up the queue.

**Response 201**
```json
{
  "id": 42,
  "email": "user@example.com",
  "referral_code": "XYZ789",
  "position": 42,
  "created_at": "2026-05-01T10:00:00Z"
}
```

Share your `referral_code` — each person who signs up with it moves you up the list.

---

## Get waitlist stats

```http
GET /waitlist/stats
```

**Response 200**
```json
{
  "total": 1247
}
```

---

## Get entry by referral code

```http
GET /waitlist/:code
```

**Response 200**
```json
{
  "id": 42,
  "email": "user@example.com",
  "referral_code": "XYZ789",
  "position": 38,
  "created_at": "2026-05-01T10:00:00Z"
}
```
