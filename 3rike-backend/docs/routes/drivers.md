# /api/drivers — Driver Profiles

All routes require `Authorization: Bearer <token>`.

A driver is a tricycle operator working toward ownership over 70 weeks.

---

## Create driver profile

```http
POST /api/drivers
Authorization: Bearer <token>
```

```json
{
  "full_name": "Emeka Okafor",
  "phone": "+2348012345678",
  "country": "Nigeria"
}
```

The `user_id` is taken from the JWT — one driver profile per user.

**Response 201**
```json
{
  "id": 1,
  "user_id": 3,
  "full_name": "Emeka Okafor",
  "phone": "+2348012345678",
  "country": "Nigeria",
  "credit_score": 0,
  "weeks_remaining": 70,
  "created_at": "2026-05-01T10:00:00Z"
}
```

`weeks_remaining` starts at 70 and decrements by 1 on each confirmed payment.

---

## Get driver by ID

```http
GET /api/drivers/:id
Authorization: Bearer <token>
```

**Response 200**
```json
{
  "id": 1,
  "user_id": 3,
  "full_name": "Emeka Okafor",
  "phone": "+2348012345678",
  "country": "Nigeria",
  "credit_score": 650,
  "weeks_remaining": 58,
  "created_at": "2026-05-01T10:00:00Z"
}
```

---

## List all drivers

```http
GET /api/drivers
Authorization: Bearer <token>
```

**Response 200** — array of driver objects.
