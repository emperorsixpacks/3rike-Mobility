package service

import (
	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/3rike12/3rike-backend/internal/repository"
	"github.com/3rike12/3rike-backend/pkg/canton"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Services struct {
	Auth     domain.AuthService
	User     domain.UserService
	Driver   domain.DriverService
	Investor domain.InvestorService
	Tricycle domain.TricycleService
	Payment  domain.PaymentService
	Loan     domain.LoanService
	Savings  domain.SavingsService
	Yield    domain.YieldService
	Waitlist domain.WaitlistService
	Fraction domain.FractionService
}

func New(db *gorm.DB, cantonClient *canton.Client, jwtSecret string, rdb *redis.Client) *Services {
	users := repository.NewUserRepo(db)
	drivers := repository.NewDriverRepo(db)
	investors := repository.NewInvestorRepo(db)
	tricycles := repository.NewTricycleRepo(db)
	payments := repository.NewPaymentRepo(db)
	loans := repository.NewLoanRepo(db)
	savings := repository.NewSavingsRepo(db)
	yields := repository.NewYieldRepo(db)
	fractions := repository.NewFractionRepo(db)
	waitlist := repository.NewWaitlistRepo(db)

	return &Services{
		Auth:     newAuthService(users, jwtSecret, rdb),
		User:     newUserService(users),
		Driver:   newDriverService(drivers),
		Investor: newInvestorService(investors, fractions),
		Tricycle: newTricycleService(tricycles, cantonClient),
		Payment:  newPaymentService(payments, drivers),
		Loan:     newLoanService(loans, drivers),
		Savings:  newSavingsService(savings),
		Yield:    newYieldService(yields, fractions),
		Waitlist: newWaitlistService(waitlist),
		Fraction: newFractionService(tricycles, investors, fractions, cantonClient),
	}
}
