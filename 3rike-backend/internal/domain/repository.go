package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uint) error
}

// InvestorRepository extended with investment listing.


type DriverRepository interface {
	Create(ctx context.Context, d *Driver) error
	FindByID(ctx context.Context, id uint) (*Driver, error)
	FindByUserID(ctx context.Context, userID uint) (*Driver, error)
	Update(ctx context.Context, d *Driver) error
	List(ctx context.Context) ([]Driver, error)
}

type InvestorRepository interface {
	Create(ctx context.Context, i *Investor) error
	FindByID(ctx context.Context, id uint) (*Investor, error)
	FindByUserID(ctx context.Context, userID uint) (*Investor, error)
	List(ctx context.Context) ([]Investor, error)
}

type TricycleRepository interface {
	Create(ctx context.Context, t *Tricycle) error
	FindByID(ctx context.Context, id uint) (*Tricycle, error)
	Update(ctx context.Context, t *Tricycle) error
	List(ctx context.Context) ([]Tricycle, error)
}

type FractionRepository interface {
	Create(ctx context.Context, f *Fraction) error
	FindByInvestorID(ctx context.Context, investorID uint) ([]Fraction, error)
	FindByTricycleID(ctx context.Context, tricycleID uint) ([]Fraction, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, p *Payment) error
	FindByDriverID(ctx context.Context, driverID uint) ([]Payment, error)
	FindByID(ctx context.Context, id uint) (*Payment, error)
	Update(ctx context.Context, p *Payment) error
}

type LoanRepository interface {
	Create(ctx context.Context, l *Loan) error
	FindByID(ctx context.Context, id uint) (*Loan, error)
	FindByDriverID(ctx context.Context, driverID uint) ([]Loan, error)
	Update(ctx context.Context, l *Loan) error
}

type SavingsRepository interface {
	FindOrCreate(ctx context.Context, driverID uint) (*Savings, error)
	FindByDriverID(ctx context.Context, driverID uint) (*Savings, error)
	Update(ctx context.Context, s *Savings) error
}

type YieldPayoutRepository interface {
	Create(ctx context.Context, y *YieldPayout) error
	FindByInvestorID(ctx context.Context, investorID uint) ([]YieldPayout, error)
}

type WaitlistRepository interface {
	Create(ctx context.Context, e *WaitlistEntry) error
	FindByEmail(ctx context.Context, email string) (*WaitlistEntry, error)
	FindByReferralCode(ctx context.Context, code string) (*WaitlistEntry, error)
	CountReferralsBy(ctx context.Context, code string) (int64, error)
	Count(ctx context.Context) (int64, error)
}
