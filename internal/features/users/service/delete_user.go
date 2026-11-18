package users_service

import (
	"context"
	"fmt"
	core_error "practice/internal/core/errors"
)

func (s *UsersService) DeleteUser(ctx context.Context, actorID, targetID int) error {
	if actorID != targetID {
		return core_error.ErrPermissionDenied
	}
	if err := s.usersRepository.DeleteUser(ctx, targetID); err != nil {
		return fmt.Errorf("delete user from repository: %w", err)
	}

	return nil
}
