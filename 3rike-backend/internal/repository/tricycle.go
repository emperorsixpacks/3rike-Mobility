package repository

import (
	"context"

	"github.com/3rike12/3rike-backend/internal/domain"
	"gorm.io/gorm"
)

type tricycleRepo struct{ db *gorm.DB }

func NewTricycleRepo(db *gorm.DB) domain.TricycleRepository { return &tricycleRepo{db} }

func (r *tricycleRepo) Create(ctx context.Context, t *domain.Tricycle) error {
	m := &Tricycle{DriverID: t.DriverID, Make: t.Make, Model: t.Model, IsEV: t.IsEV, PriceUSD: t.PriceUSD, Status: string(domain.StatusAvailable)}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	t.ID, t.CreatedAt, t.Status = m.ID, m.CreatedAt, domain.TricycleStatus(m.Status)
	return nil
}

func (r *tricycleRepo) FindByID(ctx context.Context, id uint) (*domain.Tricycle, error) {
	var m Tricycle
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomainTricycle(&m), nil
}

func (r *tricycleRepo) Update(ctx context.Context, t *domain.Tricycle) error {
	return r.db.WithContext(ctx).Save(&Tricycle{
		ID: t.ID, DriverID: t.DriverID, Make: t.Make, Model: t.Model, IsEV: t.IsEV,
		PriceUSD: t.PriceUSD, Status: string(t.Status), ContractID: t.ContractID, TotalFractions: t.TotalFractions,
	}).Error
}

func (r *tricycleRepo) List(ctx context.Context) ([]domain.Tricycle, error) {
	var ms []Tricycle
	if err := r.db.WithContext(ctx).Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Tricycle, len(ms))
	for i, m := range ms {
		out[i] = *toDomainTricycle(&m)
	}
	return out, nil
}

func toDomainTricycle(m *Tricycle) *domain.Tricycle {
	return &domain.Tricycle{ID: m.ID, DriverID: m.DriverID, Make: m.Make, Model: m.Model, IsEV: m.IsEV, PriceUSD: m.PriceUSD, Status: domain.TricycleStatus(m.Status), ContractID: m.ContractID, TotalFractions: m.TotalFractions, CreatedAt: m.CreatedAt}
}
