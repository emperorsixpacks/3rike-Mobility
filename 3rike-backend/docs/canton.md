# Canton / DAML Integration

## Overview

3riKE uses the Canton Network to tokenize tricycles and fractionalize them into investor shares. Smart contracts are written in DAML and deployed to the Canton ledger.

## Devnet setup

1. Register at the wallet UI:
   ```
   https://wallet.validator.hackcanton-01.devnet.naas.noders.services
   ```
   Use your Keycloak account (`noders-appsfactory` realm).

2. Your party ID will be shown in the wallet. It looks like:
   ```
   70af8ee8-8bc6-4a66-81a6-d375a678273a::1220195a56748e538153ecc527422256c235ff27b367483b04e161d3bbc62b1ebf32
   ```
   - Part before `::` = your Keycloak subject ID
   - Part after `::` = participant node fingerprint (same for all users on this devnet)

3. Set in `.env`:
   ```env
   CANTON_OPERATOR_PARTY=<your-party-id>
   OIDC_USERNAME=your@email.com
   OIDC_PASSWORD=your-password
   ```

## DAML contracts

Located in `../3rike-daml/daml/`:

| Contract | Purpose |
|---|---|
| `TricycleToken` | Tokenized tricycle — fractionalize, record payment, transfer ownership |
| `Fraction` | Investor share in a tricycle |
| `SavingsAccount` | Driver on-chain savings |
| `YieldPayout` | Weekly yield distribution record |

Package hash (deployed on devnet):
```
fb36f0cfcfe6c4b2a24f458f5ba06bfc697fa0584b13f44ae3d3568a294d4c19
```

## Tokenize a tricycle

```http
POST /api/tricycles/:id/tokenize
Authorization: Bearer <token>
```

- Submits a `CreateCommand` for `TricycleToken` to the Canton ledger
- `actAs` = caller's Canton party ID (from JWT)
- Stores the returned `contractId` on the tricycle record
- Status changes: `available` → `tokenized`

## Fractionalize a tricycle

```http
POST /api/tricycles/:id/fractionalize
Authorization: Bearer <token>

{ "total_fractions": 100 }
```

- Exercises the `Fractionalize` choice on the `TricycleToken` contract
- Creates `Fraction` contracts for investors
- Status changes: `tokenized` → `fractionalized`

## Stub mode

If `CANTON_URL` is empty, all Canton calls return stub data — no ledger required. Useful for local development without devnet access.

## Canton JSON Ledger API

Base URL: `https://ledger-api-json.participant.hackcanton-01.devnet.naas.noders.services`

Key endpoints used:
```
POST /v2/commands/submit-and-wait-for-transaction-tree  — submit commands
GET  /v2/packages                                        — list uploaded packages
GET  /v2/state/active-contracts                          — query active contracts
POST /v2/updates                                         — query transaction history
```

Auth: Bearer token from Keycloak (auto-fetched via OIDC password grant when `OIDC_*` vars are set).

## Integration tests

```bash
CANTON_URL=https://ledger-api-json.participant.hackcanton-01.devnet.naas.noders.services \
OIDC_TOKEN_URL=https://keycloak.naas.noders.services/realms/noders-appsfactory/protocol/openid-connect/token \
OIDC_CLIENT_ID=web-app-ui-hackcanton-01-devnet \
OIDC_USERNAME=your@email.com \
OIDC_PASSWORD=yourpassword \
CANTON_PARTY=<your-party-id> \
go test ./pkg/canton/... -v
```

All 3 tests should pass:
- `TestTokenProvider_FetchesToken` — Keycloak auth
- `TestTokenize_CreatesContractOnLedger` — creates real contract on-chain
- `TestFractionalize_ExercisesChoiceOnLedger` — full tokenize → fractionalize round trip
