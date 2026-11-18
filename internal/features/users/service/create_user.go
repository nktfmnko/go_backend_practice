package users_service

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
	core_auth "practice/internal/core/transport/auth"
)

func (s *UsersService) CreateUser(ctx context.Context, fullName string, phoneNumber *string) (domain.User, error) {
	claims, err := core_auth.MustUserFromContext(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("user unauthenticated: %w", err)
	}
	user := domain.NewUserUninitialized(claims.UserID, fullName, phoneNumber)
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	createdUser, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return createdUser, nil
}
