package users_service

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
	core_error "practice/internal/core/errors"
	core_auth "practice/internal/core/transport/auth"
)

func (s *UsersService) PatchUser(ctx context.Context, targetID int, patch domain.UserPatch) (domain.User, error) {
	claims, err := core_auth.MustUserFromContext(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("user unauthenticated: %w", err)
	}
	if claims.UserID != targetID || !claims.IsAdmin {
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
