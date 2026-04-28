package repository

import (
	"context"
	"errors"

	"github.com/3rike12/3rike-backend/internal/domain"
	"gorm.io/gorm"
)

type waitlistRepo struct{ db *gorm.DB }

func NewWaitlistRepo(db *gorm.DB) domain.WaitlistRepository { return &waitlistRepo{db} }

func (r *waitlistRepo) Create(ctx context.Context, e *domain.WaitlistEntry) error {
	m := &WaitlistEntry{
		Email:        e.Email,
		Phone:        e.Phone,
		ReferralCode: e.ReferralCode,
		ReferredBy:   e.ReferredBy,
		Position:     e.Position,
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	e.ID, e.CreatedAt = m.ID, m.CreatedAt
	return nil
}

func (r *waitlistRepo) FindByEmail(ctx context.Context, email string) (*domain.WaitlistEntry, error) {
	var m WaitlistEntry
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainWaitlist(&m), nil
}

func (r *waitlistRepo) FindByReferralCode(ctx context.Context, code string) (*domain.WaitlistEntry, error) {
	var m WaitlistEntry
	if err := r.db.WithContext(ctx).Where("referral_code = ?", code).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainWaitlist(&m), nil
}

func (r *waitlistRepo) CountReferralsBy(ctx context.Context, code string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&WaitlistEntry{}).Where("referred_by = ?", code).Count(&n).Error
	return n, err
}

func (r *waitlistRepo) Count(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&WaitlistEntry{}).Count(&n).Error
	return n, err
}

func toDomainWaitlist(m *WaitlistEntry) *domain.WaitlistEntry {
	return &domain.WaitlistEntry{
		ID:           m.ID,
		Email:        m.Email,
		Phone:        m.Phone,
		ReferralCode: m.ReferralCode,
		ReferredBy:   m.ReferredBy,
		Position:     m.Position,
		CreatedAt:    m.CreatedAt,
	}
}
