package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/3rike12/3rike-backend/internal/domain"
)

// ErrTricycleNotAvailable is returned when the tricycle isn't fractionalized
// (i.e. not yet open for fractional purchase).
var ErrTricycleNotAvailable = errors.New("tricycle_not_available")

// ErrInsufficientUnits is returned when the requested units exceed what's
// remaining for the tricycle.
var ErrInsufficientUnits = errors.New("insufficient_units")

// ErrInvalidUnits is returned when units is <= 0.
var ErrInvalidUnits = errors.New("invalid_units")

type fractionService struct {
	tricycles domain.TricycleRepository
	investors domain.InvestorRepository
	fractions domain.FractionRepository
}

func newFractionService(
	tricycles domain.TricycleRepository,
	investors domain.InvestorRepository,
	fractions domain.FractionRepository,
) domain.FractionService {
	return &fractionService{tricycles: tricycles, investors: investors, fractions: fractions}
}

// Buy is the user-facing "purchase fraction" entry point.
//
// Flow:
//  1. Validate units > 0
//  2. Look up the tricycle; reject if not fractionalized
//  3. Confirm units <= remaining units (total_fractions - already-sold)
//  4. Lazily create an Investor profile for the calling user if they don't
//     have one (lets a driver invest without a separate signup step)
//  5. Create the Fraction record
func (s *fractionService) Buy(
	ctx context.Context,
	userID uint,
	tricycleID uint,
	units int,
) (*domain.Fraction, error) {
	if units <= 0 {
		return nil, ErrInvalidUnits
	}

	t, err := s.tricycles.FindByID(ctx, tricycleID)
	if err != nil {
		return nil, fmt.Errorf("lookup tricycle: %w", err)
	}
	if t == nil {
		return nil, ErrTricycleNotAvailable
	}
	if t.Status != domain.StatusFractionalized || t.TotalFractions <= 0 {
		return nil, ErrTricycleNotAvailable
	}

	// Check remaining capacity.
	existing, err := s.fractions.FindByTricycleID(ctx, tricycleID)
	if err != nil {
		return nil, fmt.Errorf("count existing fractions: %w", err)
	}
	soldUnits := 0
	for _, f := range existing {
		soldUnits += f.Units
	}
	remaining := t.TotalFractions - soldUnits
	if units > remaining {
		return nil, ErrInsufficientUnits
	}

	// Resolve or lazily create the investor profile.
	investor, err := s.investors.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lookup investor: %w", err)
	}
	if investor == nil {
		investor = &domain.Investor{
			UserID:        userID,
			FullName:      "", // filled in via /api/investors PUT later
			WalletAddress: "",
		}
		if err := s.investors.Create(ctx, investor); err != nil {
			return nil, fmt.Errorf("create investor: %w", err)
		}
	}

	// Compute price-per-unit. Even split of price across all fractions.
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
