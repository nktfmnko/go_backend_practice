package users_redis_repository

import (
	"context"
	"practice/internal/core/domain"
)

func (r *CachedUsersRepository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	patched, err := r.next.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, err
	}

	_ = r.rdb.Del(ctx, r.getUserKey(id)).Err()

	return patched, nil
}
