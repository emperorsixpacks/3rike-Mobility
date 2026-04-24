package service

import (
	"context"
	"errors"

	"github.com/3rike12/3rike-backend/internal/domain"
)

// --- Payment ---

type paymentService struct {
	payments domain.PaymentRepository
	drivers  domain.DriverRepository
}

func newPaymentService(payments domain.PaymentRepository, drivers domain.DriverRepository) domain.PaymentService {
	return &paymentService{payments: payments, drivers: drivers}
}

func (s *paymentService) RecordPayment(ctx context.Context, p *domain.Payment) (*domain.Payment, error) {
	p.Status = domain.PaymentConfirmed
	if err := s.payments.Create(ctx, p); err != nil {
		return nil, err
	}
	// Decrement driver weeks remaining on confirmed payment.
	driver, err := s.drivers.FindByID(ctx, p.DriverID)
	if err == nil && driver.WeeksRemaining > 0 {
		driver.WeeksRemaining--
		_ = s.drivers.Update(ctx, driver)
	}
	return p, nil
}

func (s *paymentService) GetByDriver(ctx context.Context, driverID uint) ([]domain.Payment, error) {
	return s.payments.FindByDriverID(ctx, driverID)
}

// --- Loan ---

type loanService struct {
	loans   domain.LoanRepository
	drivers domain.DriverRepository
}

func newLoanService(loans domain.LoanRepository, drivers domain.DriverRepository) domain.LoanService {
	return &loanService{loans: loans, drivers: drivers}
}

func (s *loanService) Apply(ctx context.Context, l *domain.Loan) (*domain.Loan, error) {
	driver, err := s.drivers.FindByID(ctx, l.DriverID)
	if err != nil {
		return nil, errors.New("driver not found")
	}
	if driver.CreditScore < 500 {
		return nil, errors.New("credit score too low for loan eligibility")
	}
	return l, s.loans.Create(ctx, l)
}

func (s *loanService) GetByID(ctx context.Context, id uint) (*domain.Loan, error) {
	return s.loans.FindByID(ctx, id)
}

func (s *loanService) Repay(ctx context.Context, loanID uint, amountUSDC float64) (*domain.Loan, error) {
	loan, err := s.loans.FindByID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	if loan.Status != domain.LoanActive {
		return nil, errors.New("loan is not active")
	}
	loan.RemainingUSDC -= amountUSDC
	if loan.RemainingUSDC <= 0 {
		loan.RemainingUSDC = 0
		loan.Status = domain.LoanRepaid
	}
	return loan, s.loans.Update(ctx, loan)
}

// --- Savings ---

type savingsService struct{ repo domain.SavingsRepository }

func newSavingsService(repo domain.SavingsRepository) domain.SavingsService {
	return &savingsService{repo}
}

func (s *savingsService) Deposit(ctx context.Context, driverID uint, amountUSDC float64) (*domain.Savings, error) {
	if amountUSDC <= 0 {
		return nil, errors.New("deposit amount must be positive")
	}
	acc, err := s.repo.FindOrCreate(ctx, driverID)
	if err != nil {
		return nil, err
	}
	acc.BalanceUSDC += amountUSDC
	return acc, s.repo.Update(ctx, acc)
}

func (s *savingsService) GetBalance(ctx context.Context, driverID uint) (*domain.Savings, error) {
	return s.repo.FindByDriverID(ctx, driverID)
}

// --- Yield ---

type yieldService struct {
	yields    domain.YieldPayoutRepository
	fractions domain.FractionRepository
}

func newYieldService(yields domain.YieldPayoutRepository, fractions domain.FractionRepository) domain.YieldService {
	return &yieldService{yields: yields, fractions: fractions}
}

func (s *yieldService) GetByInvestor(ctx context.Context, investorID uint) ([]domain.YieldPayout, error) {
	return s.yields.FindByInvestorID(ctx, investorID)
}

// DistributeWeekly calculates and records yield payouts for all fractions.
// Yield per fraction unit = fixed 2% weekly on price per unit (placeholder rate).
func (s *yieldService) DistributeWeekly(ctx context.Context, weekNumber int) error {
	// NOTE: In production this would iterate all active tricycles and their fractions.
	// Placeholder: no-op until tricycle payment data is wired in.
	return nil
}
