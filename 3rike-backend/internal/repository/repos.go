package repository

import (
	"context"

	"github.com/3rike12/3rike-backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// --- Fraction ---

type fractionRepo struct{ db *gorm.DB }

func NewFractionRepo(db *gorm.DB) domain.FractionRepository { return &fractionRepo{db} }

func (r *fractionRepo) Create(ctx context.Context, f *domain.Fraction) error {
	m := &Fraction{TricycleID: f.TricycleID, InvestorID: f.InvestorID, Units: f.Units, PricePerUnit: f.PricePerUnit}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	f.ID, f.CreatedAt = m.ID, m.CreatedAt
	return nil
}

func (r *fractionRepo) FindByInvestorID(ctx context.Context, investorID uint) ([]domain.Fraction, error) {
	var ms []Fraction
	if err := r.db.WithContext(ctx).Where("investor_id = ?", investorID).Find(&ms).Error; err != nil {
		return nil, err
	}
	return toFractions(ms), nil
}

func (r *fractionRepo) FindByTricycleID(ctx context.Context, tricycleID uint) ([]domain.Fraction, error) {
	var ms []Fraction
	if err := r.db.WithContext(ctx).Where("tricycle_id = ?", tricycleID).Find(&ms).Error; err != nil {
		return nil, err
	}
	return toFractions(ms), nil
}

func toFractions(ms []Fraction) []domain.Fraction {
	out := make([]domain.Fraction, len(ms))
	for i, m := range ms {
		out[i] = domain.Fraction{ID: m.ID, TricycleID: m.TricycleID, InvestorID: m.InvestorID, Units: m.Units, PricePerUnit: m.PricePerUnit, CreatedAt: m.CreatedAt}
	}
	return out
}

// --- Payment ---

type paymentRepo struct{ db *gorm.DB }

func NewPaymentRepo(db *gorm.DB) domain.PaymentRepository { return &paymentRepo{db} }

func (r *paymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	m := &Payment{DriverID: p.DriverID, TricycleID: p.TricycleID, AmountLocal: p.AmountLocal, AmountUSDC: p.AmountUSDC, Currency: p.Currency, Status: string(p.Status), WeekNumber: p.WeekNumber}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	p.ID, p.CreatedAt = m.ID, m.CreatedAt
	return nil
}

func (r *paymentRepo) FindByDriverID(ctx context.Context, driverID uint) ([]domain.Payment, error) {
	var ms []Payment
	if err := r.db.WithContext(ctx).Where("driver_id = ?", driverID).Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Payment, len(ms))
	for i, m := range ms {
		out[i] = domain.Payment{ID: m.ID, DriverID: m.DriverID, TricycleID: m.TricycleID, AmountLocal: m.AmountLocal, AmountUSDC: m.AmountUSDC, Currency: m.Currency, Status: domain.PaymentStatus(m.Status), WeekNumber: m.WeekNumber, CreatedAt: m.CreatedAt}
	}
	return out, nil
}

func (r *paymentRepo) FindByID(ctx context.Context, id uint) (*domain.Payment, error) {
	var m Payment
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &domain.Payment{ID: m.ID, DriverID: m.DriverID, TricycleID: m.TricycleID, AmountLocal: m.AmountLocal, AmountUSDC: m.AmountUSDC, Currency: m.Currency, Status: domain.PaymentStatus(m.Status), WeekNumber: m.WeekNumber, CreatedAt: m.CreatedAt}, nil
}

func (r *paymentRepo) Update(ctx context.Context, p *domain.Payment) error {
	return r.db.WithContext(ctx).Save(&Payment{ID: p.ID, DriverID: p.DriverID, TricycleID: p.TricycleID, AmountLocal: p.AmountLocal, AmountUSDC: p.AmountUSDC, Currency: p.Currency, Status: string(p.Status), WeekNumber: p.WeekNumber}).Error
}

// --- Loan ---

type loanRepo struct{ db *gorm.DB }

func NewLoanRepo(db *gorm.DB) domain.LoanRepository { return &loanRepo{db} }

func (r *loanRepo) Create(ctx context.Context, l *domain.Loan) error {
	m := &Loan{DriverID: l.DriverID, PrincipalUSDC: l.PrincipalUSDC, RemainingUSDC: l.PrincipalUSDC, WeeklyRepayment: l.WeeklyRepayment, Status: string(domain.LoanActive)}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	l.ID, l.CreatedAt, l.Status, l.RemainingUSDC = m.ID, m.CreatedAt, domain.LoanActive, m.RemainingUSDC
	return nil
}

func (r *loanRepo) FindByID(ctx context.Context, id uint) (*domain.Loan, error) {
	var m Loan
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomainLoan(&m), nil
}

func (r *loanRepo) FindByDriverID(ctx context.Context, driverID uint) ([]domain.Loan, error) {
	var ms []Loan
	if err := r.db.WithContext(ctx).Where("driver_id = ?", driverID).Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Loan, len(ms))
	for i, m := range ms {
		out[i] = *toDomainLoan(&m)
	}
	return out, nil
}

func (r *loanRepo) Update(ctx context.Context, l *domain.Loan) error {
	return r.db.WithContext(ctx).Save(&Loan{ID: l.ID, DriverID: l.DriverID, PrincipalUSDC: l.PrincipalUSDC, RemainingUSDC: l.RemainingUSDC, WeeklyRepayment: l.WeeklyRepayment, Status: string(l.Status)}).Error
}

func toDomainLoan(m *Loan) *domain.Loan {
	return &domain.Loan{ID: m.ID, DriverID: m.DriverID, PrincipalUSDC: m.PrincipalUSDC, RemainingUSDC: m.RemainingUSDC, WeeklyRepayment: m.WeeklyRepayment, Status: domain.LoanStatus(m.Status), CreatedAt: m.CreatedAt}
}

// --- Savings ---

type savingsRepo struct{ db *gorm.DB }

func NewSavingsRepo(db *gorm.DB) domain.SavingsRepository { return &savingsRepo{db} }

func (r *savingsRepo) FindOrCreate(ctx context.Context, driverID uint) (*domain.Savings, error) {
	m := &Savings{DriverID: driverID}
	if err := r.db.WithContext(ctx).Where(Savings{DriverID: driverID}).Attrs(Savings{BalanceUSDC: 0}).FirstOrCreate(m).Error; err != nil {
		return nil, err
	}
	return &domain.Savings{ID: m.ID, DriverID: m.DriverID, BalanceUSDC: m.BalanceUSDC, CreatedAt: m.CreatedAt}, nil
}

func (r *savingsRepo) FindByDriverID(ctx context.Context, driverID uint) (*domain.Savings, error) {
	var m Savings
	if err := r.db.WithContext(ctx).Where("driver_id = ?", driverID).First(&m).Error; err != nil {
		return nil, err
	}
	return &domain.Savings{ID: m.ID, DriverID: m.DriverID, BalanceUSDC: m.BalanceUSDC, CreatedAt: m.CreatedAt}, nil
}

func (r *savingsRepo) Update(ctx context.Context, s *domain.Savings) error {
	return r.db.WithContext(ctx).Model(&Savings{}).Where("id = ?", s.ID).Update("balance_usdc", s.BalanceUSDC).Error
}

// --- YieldPayout ---

type yieldRepo struct{ db *gorm.DB }

func NewYieldRepo(db *gorm.DB) domain.YieldPayoutRepository { return &yieldRepo{db} }

func (r *yieldRepo) Create(ctx context.Context, y *domain.YieldPayout) error {
	m := &YieldPayout{InvestorID: y.InvestorID, FractionID: y.FractionID, AmountUSDC: y.AmountUSDC, WeekNumber: y.WeekNumber}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	y.ID, y.CreatedAt = m.ID, m.CreatedAt
	return nil
}

func (r *yieldRepo) FindByInvestorID(ctx context.Context, investorID uint) ([]domain.YieldPayout, error) {
	var ms []YieldPayout
	if err := r.db.WithContext(ctx).Where("investor_id = ?", investorID).Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.YieldPayout, len(ms))
	for i, m := range ms {
		out[i] = domain.YieldPayout{ID: m.ID, InvestorID: m.InvestorID, FractionID: m.FractionID, AmountUSDC: m.AmountUSDC, WeekNumber: m.WeekNumber, CreatedAt: m.CreatedAt}
	}
	return out, nil
}

// ensure clause import is used (upsert available if needed)
var _ = clause.OnConflict{}
