package repository

import (
	"context"

	"github.com/3rike12/3rike-backend/internal/domain"
	"gorm.io/gorm"
)

type driverRepo struct{ db *gorm.DB }

func NewDriverRepo(db *gorm.DB) domain.DriverRepository { return &driverRepo{db} }

func (r *driverRepo) Create(ctx context.Context, d *domain.Driver) error {
	m := &Driver{UserID: d.UserID, FullName: d.FullName, Phone: d.Phone, Country: d.Country, CreditScore: d.CreditScore, WeeksRemaining: 70}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	d.ID, d.CreatedAt, d.WeeksRemaining = m.ID, m.CreatedAt, m.WeeksRemaining
	return nil
}

func (r *driverRepo) FindByID(ctx context.Context, id uint) (*domain.Driver, error) {
	var m Driver
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomainDriver(&m), nil
}

func (r *driverRepo) FindByUserID(ctx context.Context, userID uint) (*domain.Driver, error) {
	var m Driver
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error; err != nil {
		return nil, err
	}
	return toDomainDriver(&m), nil
}

func (r *driverRepo) Update(ctx context.Context, d *domain.Driver) error {
	return r.db.WithContext(ctx).Save(&Driver{
		ID: d.ID, UserID: d.UserID, FullName: d.FullName, Phone: d.Phone,
		Country: d.Country, CreditScore: d.CreditScore, WeeksRemaining: d.WeeksRemaining,
	}).Error
}

func (r *driverRepo) List(ctx context.Context) ([]domain.Driver, error) {
	var ms []Driver
	if err := r.db.WithContext(ctx).Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Driver, len(ms))
	for i, m := range ms {
		out[i] = *toDomainDriver(&m)
	}
	return out, nil
}

func toDomainDriver(m *Driver) *domain.Driver {
	return &domain.Driver{ID: m.ID, UserID: m.UserID, FullName: m.FullName, Phone: m.Phone, Country: m.Country, CreditScore: m.CreditScore, WeeksRemaining: m.WeeksRemaining, CreatedAt: m.CreatedAt}
}
