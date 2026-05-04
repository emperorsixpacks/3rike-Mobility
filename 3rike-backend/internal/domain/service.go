package domain

import "context"

type AuthService interface {
	Register(ctx context.Context, email, password string, role Role) (*User, error)
	Login(ctx context.Context, email, password string) (token string, sessionID string, err error)
	Logout(ctx context.Context, sessionID string, userID uint) error
	ListSessions(ctx context.Context, userID uint) ([]Session, error)
	RevokeSession(ctx context.Context, sessionID string, userID uint) error
}

// UserService handles authenticated user self-management.
type UserService interface {
	Me(ctx context.Context, userID uint) (*User, error)
	UpdateProfile(ctx context.Context, userID uint, email string) (*User, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	DeleteAccount(ctx context.Context, userID uint) error
	LinkWallet(ctx context.Context, userID uint, cantonPartyID string) (*User, error)
}

type DriverService interface {
	Create(ctx context.Context, userID uint, d *Driver) (*Driver, error)
	GetByID(ctx context.Context, id uint) (*Driver, error)
	List(ctx context.Context) ([]Driver, error)
}

type InvestorService interface {
	Create(ctx context.Context, userID uint, i *Investor) (*Investor, error)
	GetByID(ctx context.Context, id uint) (*Investor, error)
	List(ctx context.Context) ([]Investor, error)
	ListInvestments(ctx context.Context, investorID uint) ([]Fraction, error)
}

type TricycleService interface {
	Create(ctx context.Context, t *Tricycle) (*Tricycle, error)
	GetByID(ctx context.Context, id uint) (*Tricycle, error)
	List(ctx context.Context) ([]Tricycle, error)
	Tokenize(ctx context.Context, id uint, callerParty string) (*Tricycle, error)
	Fractionalize(ctx context.Context, id uint, totalFractions int, callerParty string) (*Tricycle, error)
}

type PaymentService interface {
	RecordPayment(ctx context.Context, p *Payment) (*Payment, error)
	GetByDriver(ctx context.Context, driverID uint) ([]Payment, error)
}

type LoanService interface {
	Apply(ctx context.Context, l *Loan) (*Loan, error)
	GetByID(ctx context.Context, id uint) (*Loan, error)
	Repay(ctx context.Context, loanID uint, amountUSDC float64) (*Loan, error)
}

type SavingsService interface {
	Deposit(ctx context.Context, driverID uint, amountUSDC float64) (*Savings, error)
	GetBalance(ctx context.Context, driverID uint) (*Savings, error)
}

type YieldService interface {
	GetByInvestor(ctx context.Context, investorID uint) ([]YieldPayout, error)
	DistributeWeekly(ctx context.Context, weekNumber int) error
}

type WaitlistService interface {
	Join(ctx context.Context, email, phone string, referredBy *string) (entry *WaitlistEntry, totalCount int64, err error)
	Stats(ctx context.Context) (totalCount int64, err error)
	GetByCode(ctx context.Context, code string) (entry *WaitlistEntry, totalCount int64, referralCount int64, err error)
}

// FractionService handles user-driven purchases of tricycle fractions
// (the "buy 2 fleets / 10 shares" flow). The investor profile is lazily
// created on first purchase so callers don't need a separate signup step.
type FractionService interface {
	Buy(ctx context.Context, userID uint, tricycleID uint, units int) (*Fraction, error)
	Available(ctx context.Context, tricycleID uint) (total int, sold int, remaining int, err error)
}
