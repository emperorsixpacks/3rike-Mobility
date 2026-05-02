package repository

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint   `gorm:"primaryKey"`
	Email         string `gorm:"uniqueIndex;not null"`
	PasswordHash  string `gorm:"not null"`
	Role          string `gorm:"not null"`
	CantonPartyID string
	CreatedAt     time.Time
}

type Driver struct {
	ID             uint `gorm:"primaryKey"`
	UserID         uint `gorm:"uniqueIndex;not null"`
	FullName       string
	Phone          string
	Country        string
	CreditScore    int
	WeeksRemaining int
	CreatedAt      time.Time
}

type Investor struct {
	ID            uint `gorm:"primaryKey"`
	UserID        uint `gorm:"uniqueIndex;not null"`
	FullName      string
	WalletAddress string
	CreatedAt     time.Time
}

type Tricycle struct {
	ID             uint `gorm:"primaryKey"`
	DriverID       uint `gorm:"not null"`
	Make           string
	Model          string
	IsEV           bool
	PriceUSD       float64
	Status         string
	ContractID     string
	TotalFractions int
	CreatedAt      time.Time
}

type Fraction struct {
	ID           uint `gorm:"primaryKey"`
	TricycleID   uint `gorm:"not null"`
	InvestorID   uint `gorm:"not null"`
	Units        int
	PricePerUnit float64
	CreatedAt    time.Time
}

type Payment struct {
	ID          uint `gorm:"primaryKey"`
	DriverID    uint `gorm:"not null"`
	TricycleID  uint `gorm:"not null"`
	AmountLocal float64
	AmountUSDC  float64
	Currency    string
	Status      string
	WeekNumber  int
	CreatedAt   time.Time
}

type Loan struct {
	ID              uint `gorm:"primaryKey"`
	DriverID        uint `gorm:"not null"`
	PrincipalUSDC   float64
	RemainingUSDC   float64
	WeeklyRepayment float64
	Status          string
	CreatedAt       time.Time
}

type Savings struct {
	ID          uint `gorm:"primaryKey"`
	DriverID    uint `gorm:"uniqueIndex;not null"`
	BalanceUSDC float64
	CreatedAt   time.Time
}

type YieldPayout struct {
	ID         uint `gorm:"primaryKey"`
	InvestorID uint `gorm:"not null"`
	FractionID uint `gorm:"not null"`
	AmountUSDC float64
	WeekNumber int
	CreatedAt  time.Time
}

type WaitlistEntry struct {
	ID           uint    `gorm:"primaryKey"`
	Email        string  `gorm:"uniqueIndex;size:255;not null"`
	Phone        string  `gorm:"size:32"`
	ReferralCode string  `gorm:"uniqueIndex;size:16;not null"`
	ReferredBy   *string `gorm:"size:16;index"`
	Position     int     `gorm:"not null"`
	CreatedAt    time.Time
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{}, &Driver{}, &Investor{}, &Tricycle{},
		&Fraction{}, &Payment{}, &Loan{}, &Savings{}, &YieldPayout{},
		&WaitlistEntry{},
	)
}
