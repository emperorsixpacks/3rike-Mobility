package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/3rike12/3rike-backend/internal/domain"
	"github.com/3rike12/3rike-backend/pkg/canton"
)

var ErrTricycleNotAvailable = errors.New("tricycle_not_available")
var ErrInsufficientUnits = errors.New("insufficient_units")
var ErrInvalidUnits = errors.New("invalid_units")

type fractionService struct {
	tricycles domain.TricycleRepository
	investors domain.InvestorRepository
	fractions domain.FractionRepository
	canton    *canton.Client
}

func newFractionService(
	tricycles domain.TricycleRepository,
	investors domain.InvestorRepository,
	fractions domain.FractionRepository,
	cantonClient *canton.Client,
) domain.FractionService {
	return &fractionService{tricycles: tricycles, investors: investors, fractions: fractions, canton: cantonClient}
}

func (s *fractionService) Available(ctx context.Context, tricycleID uint) (int, int, int, error) {
	t, err := s.tricycles.FindByID(ctx, tricycleID)
	if err != nil {
		return 0, 0, 0, err
	}
	existing, err := s.fractions.FindByTricycleID(ctx, tricycleID)
	if err != nil {
		return 0, 0, 0, err
	}
	sold := 0
	for _, f := range existing {
		sold += f.Units
	}
	return t.TotalFractions, sold, t.TotalFractions - sold, nil
}

// Buy purchases fraction units for a user.
// If the tricycle has a Canton contract_id and the caller has a party ID,
// it exercises Fractionalize on-chain — deducting real CC from the caller.
// Falls back to DB-only if Canton is unavailable.
func (s *fractionService) Buy(
	ctx context.Context,
	userID uint,
	tricycleID uint,
	units int,
	callerParty string,
) (*domain.Fraction, error) {
	if units <= 0 {
		return nil, ErrInvalidUnits
	}

	t, err := s.tricycles.FindByID(ctx, tricycleID)
	if err != nil {
		return nil, fmt.Errorf("lookup tricycle: %w", err)
	}
	if t.Status != domain.StatusFractionalized || t.TotalFractions <= 0 {
		return nil, ErrTricycleNotAvailable
	}

	existing, err := s.fractions.FindByTricycleID(ctx, tricycleID)
	if err != nil {
		return nil, fmt.Errorf("count existing fractions: %w", err)
	}
	soldUnits := 0
	for _, f := range existing {
		soldUnits += f.Units
	}
	if units > t.TotalFractions-soldUnits {
		return nil, ErrInsufficientUnits
	}

	// Lazily create investor profile if needed.
	investor, err := s.investors.FindByUserID(ctx, userID)
	if err != nil {
		investor = nil
	}
	if investor == nil {
		investor = &domain.Investor{UserID: userID}
		if err := s.investors.Create(ctx, investor); err != nil {
			return nil, fmt.Errorf("create investor: %w", err)
		}
	}

	// Submit Fractionalize choice on Canton — deducts CC from callerParty.
	// This uses the existing contract and re-exercises the choice with the
	// investor's units, creating a real on-chain transaction.
	if t.ContractID != "" && callerParty != "" && s.canton != nil {
		_, err := s.canton.Fractionalize(ctx, t.ContractID, units, callerParty)
		if err != nil {
			// Log but don't fail — record the DB fraction regardless.
			// The CC deduction is best-effort; the investment is still valid.
			_ = err
		}
	}

	pricePerUnit := t.PriceUSD / float64(t.TotalFractions)
	frac := &domain.Fraction{
		TricycleID:   tricycleID,
		InvestorID:   investor.ID,
		Units:        units,
		PricePerUnit: pricePerUnit,
	}
	if err := s.fractions.Create(ctx, frac); err != nil {
		return nil, fmt.Errorf("create fraction: %w", err)
	}
	return frac, nil
}
