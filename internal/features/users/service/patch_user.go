package users_service

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
	core_error "practice/internal/core/errors"
)

func (s *UsersService) PatchUser(ctx context.Context, actorID, targetID int, patch domain.UserPatch) (domain.User, error) {
	if actorID != targetID {
		return domain.User{}, core_error.ErrPermissionDenied
	}
	user, err := s.usersRepository.GetUser(ctx, targetID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	patchUser, err := s.usersRepository.PatchUser(ctx, targetID, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	return patchUser, nil
}
