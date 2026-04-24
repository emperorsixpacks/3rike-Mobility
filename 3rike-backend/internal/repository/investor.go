package repository

import (
	"context"

	"github.com/3rike12/3rike-backend/internal/domain"
	"gorm.io/gorm"
)

type investorRepo struct{ db *gorm.DB }

func NewInvestorRepo(db *gorm.DB) domain.InvestorRepository { return &investorRepo{db} }

func (r *investorRepo) Create(ctx context.Context, i *domain.Investor) error {
	m := &Investor{UserID: i.UserID, FullName: i.FullName, WalletAddress: i.WalletAddress}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	i.ID, i.CreatedAt = m.ID, m.CreatedAt
	return nil
}

func (r *investorRepo) FindByID(ctx context.Context, id uint) (*domain.Investor, error) {
	var m Investor
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomainInvestor(&m), nil
}

func (r *investorRepo) FindByUserID(ctx context.Context, userID uint) (*domain.Investor, error) {
	var m Investor
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error; err != nil {
		return nil, err
	}
	return toDomainInvestor(&m), nil
}

func (r *investorRepo) List(ctx context.Context) ([]domain.Investor, error) {
	var ms []Investor
	if err := r.db.WithContext(ctx).Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Investor, len(ms))
	for i, m := range ms {
		out[i] = *toDomainInvestor(&m)
	}
	return out, nil
}

func toDomainInvestor(m *Investor) *domain.Investor {
	return &domain.Investor{ID: m.ID, UserID: m.UserID, FullName: m.FullName, WalletAddress: m.WalletAddress, CreatedAt: m.CreatedAt}
}
