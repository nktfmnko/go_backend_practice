package users_redis_repository

import "context"

func (r *CachedUsersRepository) DeleteUser(ctx context.Context, id int) error {
	key := r.getUserKey(id)
	if err := r.next.DeleteUser(ctx, id); err != nil {
		return err
	}

	_ = r.rdb.Del(ctx, key).Err()
	return nil
}
