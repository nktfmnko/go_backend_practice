package tasks_redis_repository

import (
	"context"
	"log"
)

func (r *CachedTasksRepository) DeleteTask(ctx context.Context, id int) error {
	if err := r.next.DeleteTask(ctx, id); err != nil {
		return err
	}

	if err := r.rdb.Del(ctx, r.getTaskKey(id)).Err(); err != nil {
		log.Printf("ERROR: failed to delete task %d from cache: %v", id, err)
	}
	return nil
}
