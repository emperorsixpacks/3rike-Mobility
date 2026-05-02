# API Reference

Base URL: `http://localhost:8080`

All protected routes require: `Authorization: Bearer <token>`

Interactive docs: `GET /docs/` (open in browser)

---

## Auth (public)

| Method | Path | Description |
|---|---|---|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | Login, get JWT + session |

## User (protected)

| Method | Path | Description |
|---|---|---|
| GET | `/auth/me` | Current user profile |
| PUT | `/auth/profile` | Update email |
| PUT | `/auth/password` | Change password |
| DELETE | `/auth/account` | Delete account |
| POST | `/auth/logout` | Logout current session |
| GET | `/auth/sessions` | List active sessions |
| DELETE | `/auth/sessions/:sessionID` | Revoke a session |
| PUT | `/auth/wallet` | Link Canton party ID |
| GET | `/auth/wallet/balance` | Get CC balance from Canton |

## Drivers (protected)

| Method | Path | Description |
|---|---|---|
| POST | `/api/drivers` | Create driver profile |
| GET | `/api/drivers` | List all drivers |
| GET | `/api/drivers/:id` | Get driver by ID |

## Investors (protected)

| Method | Path | Description |
|---|---|---|
| POST | `/api/investors` | Create investor profile |
| GET | `/api/investors` | List all investors |
| GET | `/api/investors/:id` | Get investor by ID |
| GET | `/api/investors/:id/investments` | List investor's fractions |

## Tricycles (protected)

| Method | Path | Description |
|---|---|---|
| POST | `/api/tricycles` | Register a tricycle |
| GET | `/api/tricycles` | List all tricycles |
| GET | `/api/tricycles/:id` | Get tricycle by ID |
| POST | `/api/tricycles/:id/tokenize` | Tokenize on Canton ledger |
| POST | `/api/tricycles/:id/fractionalize` | Fractionalize into investor shares |

## Payments (protected)

| Method | Path | Description |
|---|---|---|
| POST | `/api/payments` | Record a driver weekly payment |
| GET | `/api/payments/driver/:driverID` | Get payments for a driver |

## Loans (protected)

| Method | Path | Description |
|---|---|---|
| POST | `/api/loans` | Apply for a loan (credit score ≥ 500) |
| GET | `/api/loans/:id` | Get loan by ID |
| PUT | `/api/loans/:id/repay` | Make a repayment |

## Savings (protected)

| Method | Path | Description |
|---|---|---|
| POST | `/api/savings/deposit` | Deposit USDC into savings |
| GET | `/api/savings/:driverID/balance` | Get savings balance |

## Yield (protected)

| Method | Path | Description |
|---|---|---|
| GET | `/api/yield/investor/:investorID` | Get yield payouts for investor |

## Waitlist (public)

| Method | Path | Description |
|---|---|---|
| POST | `/waitlist/join` | Join the waitlist |
| GET | `/waitlist/stats` | Waitlist stats |
| GET | `/waitlist/:code` | Get entry by referral code |

---

## Common request/response patterns

### Error response
```json
{ "error": "description of what went wrong" }
```

### Pagination
Not yet implemented — all list endpoints return full results.

### Caching
All `GET` endpoints on `/api/*` are Redis cached:
- Drivers, investors, tricycles, loans: 5 min TTL
- Payments, savings balance: 2 min TTL
- Yield: 5 min TTL
