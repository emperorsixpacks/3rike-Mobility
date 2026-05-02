# 3riKE Backend

Real-world asset platform — tricycle ownership financing, tokenization, and investor yield on the Canton Network.

## Quick Links

- [Getting Started](docs/getting-started.md)
- [Authentication & Sessions](docs/auth.md)
- [Canton / DAML Integration](docs/canton.md)
- [API Reference](docs/api.md)
- [Architecture](docs/architecture.md)

## Stack

| Layer | Technology |
|---|---|
| HTTP | Go + Fiber v2 |
| Database | PostgreSQL (GORM) |
| Cache | Redis (session store + cache-aside) |
| Smart Contracts | DAML on Canton Network |
| Auth | JWT + Redis sessions (max 3 concurrent) |
| Docs | Swagger at `/docs/` |

## Run locally

```bash
cp .env.example .env
# fill in your values

docker-compose up -d        # starts postgres + redis
go run .                    # starts API on :8080
```

Swagger UI: http://localhost:8080/docs/
