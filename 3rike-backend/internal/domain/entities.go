// Package domain contains all business entities and repository/service interfaces.
// No external dependencies — only the Go standard library.
package domain

import "time"

type Role string

const (
	RoleDriver   Role = "driver"
	RoleInvestor Role = "investor"
	RoleAdmin    Role = "admin"
)

// User is the shared auth entity.
type User struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// Driver profile linked to a User.
type Driver struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	FullName        string    `json:"full_name"`
	Phone           string    `json:"phone"`
	Country         string    `json:"country"`
	CreditScore     int       `json:"credit_score"`
	WeeksRemaining  int       `json:"weeks_remaining"` // out of 70
	CreatedAt       time.Time `json:"created_at"`
}

// Investor profile linked to a User.
type Investor struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	FullName      string    `json:"full_name"`
	WalletAddress string    `json:"wallet_address"`
	CreatedAt     time.Time `json:"created_at"`
}

type TricycleStatus string

const (
	StatusAvailable    TricycleStatus = "available"
	StatusFinancing    TricycleStatus = "financing"
	StatusOwned        TricycleStatus = "owned"
	StatusTokenized    TricycleStatus = "tokenized"
	StatusFractionalized TricycleStatus = "fractionalized"
)

// Tricycle is the real-world asset being financed and tokenized.
type Tricycle struct {
	ID              uint           `json:"id"`
	DriverID        uint           `json:"driver_id"`
	Make            string         `json:"make"`
	Model           string         `json:"model"`
	IsEV            bool           `json:"is_ev"`
	PriceUSD        float64        `json:"price_usd"`
	Status          TricycleStatus `json:"status"`
	ContractID      string         `json:"contract_id"`       // Canton ledger contract ID
	TotalFractions  int            `json:"total_fractions"`
	CreatedAt       time.Time      `json:"created_at"`
}

// Fraction represents an investor's share in a Tricycle.
type Fraction struct {
	ID          uint      `json:"id"`
	TricycleID  uint      `json:"tricycle_id"`
	InvestorID  uint      `json:"investor_id"`
	Units       int       `json:"units"`
	PricePerUnit float64  `json:"price_per_unit"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentConfirmed PaymentStatus = "confirmed"
	PaymentFailed    PaymentStatus = "failed"
)

// Payment is a driver's weekly repayment.
type Payment struct {
	ID           uint          `json:"id"`
	DriverID     uint          `json:"driver_id"`
	TricycleID   uint          `json:"tricycle_id"`
	AmountLocal  float64       `json:"amount_local"`
	AmountUSDC   float64       `json:"amount_usdc"`
	Currency     string        `json:"currency"`
	Status       PaymentStatus `json:"status"`
	WeekNumber   int           `json:"week_number"`
	CreatedAt    time.Time     `json:"created_at"`
}

type LoanStatus string

const (
	LoanActive    LoanStatus = "active"
	LoanRepaid    LoanStatus = "repaid"
	LoanDefaulted LoanStatus = "defaulted"
)

// Loan is a credit facility for drivers with good repayment history.
type Loan struct {
	ID              uint       `json:"id"`
	DriverID        uint       `json:"driver_id"`
	PrincipalUSDC   float64    `json:"principal_usdc"`
	RemainingUSDC   float64    `json:"remaining_usdc"`
	WeeklyRepayment float64    `json:"weekly_repayment"`
	Status          LoanStatus `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
}

// Savings is a driver's savings account on the platform.
type Savings struct {
	ID        uint      `json:"id"`
	DriverID  uint      `json:"driver_id"`
	BalanceUSDC float64 `json:"balance_usdc"`
	CreatedAt time.Time `json:"created_at"`
}

// Session represents an active user session stored in Redis.
type Session struct {
	ID        string    `json:"id"`
	UserID    uint      `json:"user_id"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// YieldPayout is a weekly yield distribution to an investor.
type YieldPayout struct {
	ID         uint      `json:"id"`
	InvestorID uint      `json:"investor_id"`
	FractionID uint      `json:"fraction_id"`
	AmountUSDC float64   `json:"amount_usdc"`
	WeekNumber int       `json:"week_number"`
	CreatedAt  time.Time `json:"created_at"`
}
