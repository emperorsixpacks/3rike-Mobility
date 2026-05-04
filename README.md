# 3riKE Mobility

Real-world asset platform for tricycle ownership financing, tokenization, and investor yield on the Canton Network.

Drivers access affordable financing to own tricycles. Investors earn yield by holding fractional ownership tokens. All asset lifecycle events are recorded on-chain via DAML smart contracts.

3riKE empowers tricycle drivers in Accra to move toward ownership through signed agreements and weekly progress tracking on the Canton Network, creating a secure, verifiable record that reduces disputes and builds trust with drivers, owners, and investors.

## Monorepo structure

| Directory | Description |
|---|---|
| `3rike-backend/` | Go + Fiber REST API — auth, financing, Canton integration |
| `3rike-frontend/` | React + Vite web app — driver and investor dashboards |
| `3rike-daml/` | DAML smart contracts — TricycleToken, Savings, YieldPayout |

## How it works

1. A driver applies for tricycle financing through the platform.
2. The tricycle is tokenized on the Canton Network as a `TricycleToken` contract.
3. Investors purchase `Fraction` tokens representing fractional ownership.
4. Drivers make weekly repayments; yield is distributed to investors via `YieldPayout` contracts.
5. On full repayment, ownership transfers to the driver.

## Quick start

### Backend

```bash
cd 3rike-backend
cp .env.example .env        # fill in your values
docker-compose up -d        # start Postgres + Redis
go run .                    # API on :8080
```

Swagger UI: http://localhost:8080/docs/

### Frontend

```bash
cd 3rike-frontend
cp .env.example .env
npm install
npm run dev                 # dev server on :5173
```

### DAML contracts

```bash
cd 3rike-daml
# requires DAML SDK — curl -sSL https://get.daml.com/ | sh -s 3.1.0
daml build                  # outputs .daml/dist/3rike-contracts-1.0.0.dar
```

## Tech stack

| Layer | Technology |
|---|---|
| API | Go, Fiber v2 |
| Database | PostgreSQL (GORM) |
| Cache / Sessions | Redis |
| Smart Contracts | DAML on Canton Network |
| Auth | JWT + Redis sessions (max 3 concurrent) |
| Frontend | React 19, TypeScript, Vite, Tailwind CSS |
| i18n | i18next |

## Documentation

- [Backend README](3rike-backend/README.md)
- [Architecture](3rike-backend/docs/architecture.md)
- [Authentication](3rike-backend/docs/auth.md)
- [Canton / DAML Integration](3rike-backend/docs/canton.md)
- [API Reference](3rike-backend/docs/api.md)
- [DAML Contracts README](3rike-daml/README.md)
