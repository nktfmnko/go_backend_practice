package tasks_redis_repository

import "context"

func (r *CachedTasksRepository) DeleteTask(ctx context.Context, id int) error {
	if err := r.next.DeleteTask(ctx, id); err != nil {
		return err
	}

	_ = r.rdb.Del(ctx, r.getTaskKey(id)).Err()
	return nil
}
