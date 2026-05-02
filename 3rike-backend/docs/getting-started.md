# Getting Started

## Prerequisites

- Go 1.23+
- Docker (for Postgres + Redis)
- A Canton devnet account (see [Canton setup](canton.md))

## 1. Clone and configure

```bash
git clone <repo>
cd 3rike-backend
cp .env.example .env
```

Edit `.env`:

```env
APP_PORT=8080
ENV=development

# Postgres
DATABASE_URL=postgres://3rike:3rike@localhost:5432/3rike?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# JWT — change in production
JWT_SECRET=your-secret-here

# Canton (see docs/canton.md)
CANTON_URL=https://ledger-api-json.participant.hackcanton-01.devnet.naas.noders.services
CANTON_VALIDATOR_URL=https://wallet.validator.hackcanton-01.devnet.naas.noders.services
CANTON_OPERATOR_PARTY=<your-party-id>

# Keycloak OIDC — for Canton token auto-fetch
OIDC_TOKEN_URL=https://keycloak.naas.noders.services/realms/noders-appsfactory/protocol/openid-connect/token
OIDC_CLIENT_ID=web-app-ui-hackcanton-01-devnet
OIDC_USERNAME=your@email.com
OIDC_PASSWORD=your-password
```

## 2. Start infrastructure

```bash
docker-compose up -d
```

This starts:
- PostgreSQL on `localhost:5432`
- Redis on `localhost:6379`

## 3. Run the API

```bash
go run .
```

Or build and run:

```bash
go build -ldflags '-s -w' -o app .
./app
```

## 4. Verify

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

Swagger UI: http://localhost:8080/docs/

## 5. Run tests

```bash
# Unit/integration tests (requires DATABASE_URL)
go test ./...

# Canton integration tests (requires devnet credentials)
CANTON_URL=... OIDC_USERNAME=... OIDC_PASSWORD=... CANTON_PARTY=... \
  go test ./pkg/canton/... -v
```
