package users_service

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
)

func (s *UsersService) CreateUser(ctx context.Context, id int, fullName string, phoneNumber *string) (domain.User, error) {
	user := domain.NewUserUninitialized(id, fullName, phoneNumber)
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	createdUser, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return createdUser, nil
}
