package service

import (
	"context"
	"errors"

	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/3rike12/3rike-backend/pkg/canton"
)

type tricycleService struct {
	repo   domain.TricycleRepository
	canton *canton.Client
}

func newTricycleService(repo domain.TricycleRepository, c *canton.Client) domain.TricycleService {
	return &tricycleService{repo: repo, canton: c}
}

func (s *tricycleService) Create(ctx context.Context, t *domain.Tricycle) (*domain.Tricycle, error) {
	return t, s.repo.Create(ctx, t)
}

func (s *tricycleService) GetByID(ctx context.Context, id uint) (*domain.Tricycle, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *tricycleService) List(ctx context.Context) ([]domain.Tricycle, error) {
	return s.repo.List(ctx)
}

func (s *tricycleService) Tokenize(ctx context.Context, id uint, callerParty string) (*domain.Tricycle, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t.Status != domain.StatusAvailable && t.Status != domain.StatusFinancing {
		return nil, errors.New("tricycle cannot be tokenized in current status")
	}
	party := callerParty
	if party == "" {
		party = "operator"
	}
	res, err := s.canton.Tokenize(ctx, t.ID, party)
	if err != nil {
		return nil, err
	}
	t.ContractID = res.ContractID
	t.Status = domain.StatusTokenized
	return t, s.repo.Update(ctx, t)
}

func (s *tricycleService) Fractionalize(ctx context.Context, id uint, totalFractions int, callerParty string) (*domain.Tricycle, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t.Status != domain.StatusTokenized {
		return nil, errors.New("tricycle must be tokenized before fractionalization")
	}
	if totalFractions < 1 {
		return nil, errors.New("totalFractions must be >= 1")
	}
	party := callerParty
	if party == "" {
		party = "operator"
	}
	res, err := s.canton.Fractionalize(ctx, t.ContractID, totalFractions, party)
	if err != nil {
		return nil, err
	}
	t.ContractID = res.ContractID
	t.TotalFractions = res.TotalFractions
	t.Status = domain.StatusFractionalized
	return t, s.repo.Update(ctx, t)
}
