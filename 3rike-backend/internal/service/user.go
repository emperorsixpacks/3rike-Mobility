package service

import (
	"context"
	"errors"

	"github.com/3rike12/3rike-backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	users domain.UserRepository
}

func newUserService(users domain.UserRepository) domain.UserService {
	return &userService{users}
}

func (s *userService) Me(ctx context.Context, userID uint) (*domain.User, error) {
	return s.users.FindByID(ctx, userID)
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint, email string) (*domain.User, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	u.Email = email
	return u, s.users.Update(ctx, u)
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("current password is incorrect")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return s.users.Update(ctx, u)
}

func (s *userService) DeleteAccount(ctx context.Context, userID uint) error {
	return s.users.Delete(ctx, userID)
}
