package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/3rike12/3rike-backend/internal/domain"
)

// ErrInvalidEmail is returned when the supplied email fails parsing.
var ErrInvalidEmail = errors.New("invalid_email")

// ErrInvalidPhone is returned when a non-empty phone fails the E.164-ish check.
var ErrInvalidPhone = errors.New("invalid_phone")

// ErrInvalidReferrer is returned when referredBy doesn't resolve to an existing entry.
var ErrInvalidReferrer = errors.New("invalid_referrer")

// ErrReferralCollision indicates we couldn't generate a unique referral code
// after several attempts. Should be exceedingly rare given the keyspace.
var ErrReferralCollision = errors.New("referral_collision")

// referralCharset excludes visually ambiguous chars (0/O, 1/I/L) so codes
// shared verbally or hand-written are less error-prone.
const referralCharset = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
const referralLength = 6
const referralMaxRetries = 5

// phoneRegex is a permissive E.164 sanity check: an optional leading "+"
// followed by 8-15 digits.
var phoneRegex = regexp.MustCompile(`^\+?\d{8,15}$`)

type waitlistService struct {
	repo domain.WaitlistRepository
}

func newWaitlistService(repo domain.WaitlistRepository) domain.WaitlistService {
	return &waitlistService{repo: repo}
}

func (s *waitlistService) Join(ctx context.Context, email, phone string, referredBy *string) (*domain.WaitlistEntry, int64, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, 0, ErrInvalidEmail
	}

	phone = strings.TrimSpace(phone)
	if phone != "" && !phoneRegex.MatchString(strings.ReplaceAll(phone, " ", "")) {
		return nil, 0, ErrInvalidPhone
	}

	// Idempotent on email — re-submitting returns the existing entry.
	if existing, err := s.repo.FindByEmail(ctx, email); err != nil {
		return nil, 0, err
	} else if existing != nil {
		total, err := s.repo.Count(ctx)
		if err != nil {
			return nil, 0, err
		}
		return existing, total, nil
	}

	// Validate referredBy if supplied.
	if referredBy != nil && *referredBy != "" {
		ref, err := s.repo.FindByReferralCode(ctx, *referredBy)
		if err != nil {
			return nil, 0, err
		}
		if ref == nil {
			return nil, 0, ErrInvalidReferrer
		}
	} else {
		referredBy = nil
	}

	// Generate a unique referral code, retrying on collision.
	code, err := s.generateUniqueCode(ctx)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	entry := &domain.WaitlistEntry{
		Email:        email,
		Phone:        phone,
		ReferralCode: code,
		ReferredBy:   referredBy,
		Position:     int(count) + 1,
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, 0, err
	}

	return entry, count + 1, nil
}

func (s *waitlistService) Stats(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func (s *waitlistService) GetByCode(ctx context.Context, code string) (*domain.WaitlistEntry, int64, int64, error) {
	entry, err := s.repo.FindByReferralCode(ctx, code)
	if err != nil {
		return nil, 0, 0, err
	}
	if entry == nil {
		return nil, 0, 0, fmt.Errorf("not found")
	}
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, 0, err
	}
	refs, err := s.repo.CountReferralsBy(ctx, code)
	if err != nil {
		return nil, 0, 0, err
	}
	return entry, total, refs, nil
}

func (s *waitlistService) generateUniqueCode(ctx context.Context) (string, error) {
	for i := 0; i < referralMaxRetries; i++ {
		code, err := randomCode(referralLength)
		if err != nil {
			return "", err
		}
		existing, err := s.repo.FindByReferralCode(ctx, code)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return code, nil
		}
	}
	return "", ErrReferralCollision
}

func randomCode(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	out := make([]byte, n)
	for i, x := range b {
		out[i] = referralCharset[int(x)%len(referralCharset)]
	}
	return string(out), nil
}
