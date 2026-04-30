package handler

import "github.com/3rike12/3rike-backend/internal/service"

// Handlers is the top-level handler container.
type Handlers struct {
	Auth     *AuthHandler
	User     *UserHandler
	Driver   *DriverHandler
	Investor *InvestorHandler
	Tricycle *TricycleHandler
	Payment  *PaymentHandler
	Loan     *LoanHandler
	Savings  *SavingsHandler
	Yield    *YieldHandler
	Waitlist *WaitlistHandler
	Fraction *FractionHandler
}

func New(svc *service.Services) *Handlers {
	return &Handlers{
		Auth:     &AuthHandler{svc: svc.Auth},
		User:     &UserHandler{svc: svc.User},
		Driver:   &DriverHandler{svc: svc.Driver},
		Investor: &InvestorHandler{svc: svc.Investor},
		Tricycle: &TricycleHandler{svc: svc.Tricycle},
		Payment:  &PaymentHandler{svc: svc.Payment},
		Loan:     &LoanHandler{svc: svc.Loan},
		Savings:  &SavingsHandler{svc: svc.Savings},
		Yield:    &YieldHandler{svc: svc.Yield},
		Waitlist: &WaitlistHandler{svc: svc.Waitlist},
		Fraction: &FractionHandler{svc: svc.Fraction},
	}
}
