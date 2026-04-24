package service

import (
	"context"

	"github.com/3rike12/3rike-backend/internal/domain"
)

type driverService struct{ repo domain.DriverRepository }

func newDriverService(repo domain.DriverRepository) domain.DriverService { return &driverService{repo} }

func (s *driverService) Create(ctx context.Context, userID uint, d *domain.Driver) (*domain.Driver, error) {
	d.UserID = userID
	return d, s.repo.Create(ctx, d)
}

func (s *driverService) GetByID(ctx context.Context, id uint) (*domain.Driver, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *driverService) List(ctx context.Context) ([]domain.Driver, error) {
	return s.repo.List(ctx)
}

// ---

type investorService struct {
	repo      domain.InvestorRepository
	fractions domain.FractionRepository
}

func newInvestorService(repo domain.InvestorRepository, fractions domain.FractionRepository) domain.InvestorService {
	return &investorService{repo: repo, fractions: fractions}
}

func (s *investorService) Create(ctx context.Context, userID uint, i *domain.Investor) (*domain.Investor, error) {
	i.UserID = userID
	return i, s.repo.Create(ctx, i)
}

func (s *investorService) GetByID(ctx context.Context, id uint) (*domain.Investor, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *investorService) List(ctx context.Context) ([]domain.Investor, error) {
	return s.repo.List(ctx)
}

func (s *investorService) ListInvestments(ctx context.Context, investorID uint) ([]domain.Fraction, error) {
	return s.fractions.FindByInvestorID(ctx, investorID)
}
