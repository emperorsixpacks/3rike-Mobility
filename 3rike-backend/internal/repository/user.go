package repository

import (
	"context"

	"github.com/3rike12/3rike-backend/internal/domain"
	"gorm.io/gorm"
)

type userRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) domain.UserRepository { return &userRepo{db} }

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	m := &User{Email: u.Email, PasswordHash: u.PasswordHash, Role: string(u.Role)}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	u.ID, u.CreatedAt = m.ID, m.CreatedAt
	return nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return toDomainUser(&m), nil
}

func (r *userRepo) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var m User
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomainUser(&m), nil
}

func (r *userRepo) Update(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", u.ID).
		Updates(map[string]any{"email": u.Email, "password_hash": u.PasswordHash}).Error
}

func (r *userRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&User{}, id).Error
}

func toDomainUser(m *User) *domain.User {
	return &domain.User{ID: m.ID, Email: m.Email, PasswordHash: m.PasswordHash, Role: domain.Role(m.Role), CreatedAt: m.CreatedAt}
}
