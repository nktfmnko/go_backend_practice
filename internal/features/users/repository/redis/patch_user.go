package users_redis_repository

import (
	"context"
	"log"
	"practice/internal/core/domain"
)

func (r *CachedUsersRepository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	patched, err := r.next.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, err
	}

	if err := r.rdb.Del(ctx, r.getUserKey(id)).Err(); err != nil {
		log.Printf("ERROR: failed to delete user %d from cache: %v", id, err)
	}

	return patched, nil
}
