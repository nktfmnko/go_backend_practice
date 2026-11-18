package users_service

import (
	"context"
	"fmt"
	core_error "practice/internal/core/errors"
	core_auth "practice/internal/core/transport/auth"
)

func (s *UsersService) DeleteUser(ctx context.Context, targetID int) error {
	claims, err := core_auth.MustUserFromContext(ctx)
	if err != nil {
		return fmt.Errorf("user unauthenticated: %w", err)
	}
	if claims.UserID != targetID || !claims.IsAdmin {
		return core_error.ErrPermissionDenied
	}
	if err := s.usersRepository.DeleteUser(ctx, targetID); err != nil {
		return fmt.Errorf("delete user from repository: %w", err)
	}

	return nil
}
