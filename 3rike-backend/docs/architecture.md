# Architecture

## Clean Architecture layers

```
internal/
├── domain/          ← entities + interfaces (zero external deps)
├── repository/      ← GORM implementations (only knows domain + DB)
├── service/         ← business logic (only knows domain interfaces)
└── handler/         ← HTTP (only knows service interfaces)

pkg/
├── canton/          ← Canton JSON Ledger API client
├── middleware/       ← JWT auth + Redis cache-aside
└── testutil/        ← Postgres test bootstrap
```

**Dependency rule:** each layer only imports the layer directly below it. `handler` never imports `repository`. `service` never imports `fiber`.

## Request flow

```
HTTP Request
  → Fiber router
  → middleware.Auth (validate JWT, check Redis session)
  → middleware.Cache (return cached response if GET)
  → handler (parse request, call service)
  → service (business logic, call repository + canton)
  → repository (GORM query)
  → Postgres

Canton commands (tokenize/fractionalize):
  → service/tricycle.go
  → pkg/canton/canton.go
  → Canton JSON Ledger API
  → DAML contract on ledger
```

## Session flow

```
Login
  → bcrypt verify password
  → create session ID
  → store session:{id} in Redis (TTL 72h)
  → push to user_sessions:{userID} list
  → if list > 3, evict oldest session
  → return JWT containing session_id

Every request
  → validate JWT signature
  → check session:{id} exists in Redis
  → if missing → 401 (logged out or evicted)

Logout
  → delete session:{id} from Redis
  → remove from user_sessions list
  → token immediately invalid
```

## Canton party model

```
Operator party  = 3riKE platform (signs all contracts)
Driver party    = individual driver (observer on their tricycle)
Investor party  = individual investor (observer on their fractions)

Contract signatories: operator
actAs on commands:    caller's party (from JWT canton_party_id claim)
```

Users link their Canton party via `PUT /auth/wallet`. Until linked, commands fall back to operator party (stub/demo mode).

## Caching strategy

Redis cache-aside on all GET routes:
- Key: `cache:{path}?{querystring}`
- On hit: return cached JSON directly
- On miss: execute handler, cache 200 responses
- TTL: 2-5 minutes depending on resource
- Graceful degradation: if Redis is down, requests pass through uncached

## Database schema

```
users           — email, password_hash, role, canton_party_id
drivers         — user_id, full_name, phone, country, credit_score, weeks_remaining
investors       — user_id, full_name, wallet_address
tricycles       — driver_id, make, model, is_ev, price_usd, status, contract_id, total_fractions
fractions       — tricycle_id, investor_id, units, price_per_unit
payments        — driver_id, tricycle_id, amount_local, amount_usdc, currency, status, week_number
loans           — driver_id, principal_usdc, remaining_usdc, weekly_repayment, status
savings         — driver_id, balance_usdc
yield_payouts   — investor_id, fraction_id, amount_usdc, week_number
waitlist_entries — email, phone, referral_code, referred_by, position
```
