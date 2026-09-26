package users_redis_repository

import (
	"context"
	"log"
)

func (r *CachedUsersRepository) DeleteUser(ctx context.Context, id int) error {
	key := r.getUserKey(id)
	if err := r.next.DeleteUser(ctx, id); err != nil {
		return err
	}

	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		log.Printf("ERROR: failed to delete user %d from cache: %v", id, err)
	}
	return nil
}
