# 3riKE DAML Contracts

Smart contracts for the 3riKE platform on the Canton Network.

## Contracts

| Contract | Purpose |
|---|---|
| `TricycleToken` | Tokenized tricycle asset — fractionalize, record payments, transfer ownership |
| `Fraction` | Investor share in a tricycle |
| `SavingsAccount` | Driver on-chain savings with deposit/withdraw/interest |
| `YieldPayout` | Weekly yield distribution record for investors |

## Build

```bash
# Install DAML SDK first
curl -sSL https://get.daml.com/ | sh -s 3.1.0

# Build .dar archive
daml build

# Output: .daml/dist/3rike-contracts-1.0.0.dar
```

## Deploy to Canton node

```bash
# Upload .dar to your Canton participant node
curl -X POST http://localhost:7575/v2/packages \
  -H "Authorization: Bearer $CANTON_TOKEN" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @.daml/dist/3rike-contracts-1.0.0.dar
```

## Connecting to the Go backend

Set in `3rike-backend/.env`:
```
CANTON_URL=http://localhost:7575
CANTON_TOKEN=<your-participant-token>
```

The `templateId` used in the Go backend maps to:
```
3rike-contracts:TricycleToken:TricycleToken
3rike-contracts:Savings:SavingsAccount
3rike-contracts:YieldPayout:YieldPayout
```

## Party model

| Party | Role |
|---|---|
| `operator` | 3riKE platform — signs all contracts |
| `driver` | Observer on their tricycle token |
| `investor` | Observer on their fractions and yield payouts |
