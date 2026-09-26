package tasks_redis_repository

import (
	"context"
	"log"
	"practice/internal/core/domain"
)

func (r *CachedTasksRepository) PatchTask(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
	patched, err := r.next.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, err
	}

	if err := r.rdb.Del(ctx, r.getTaskKey(id)).Err(); err != nil {
		log.Printf("ERROR: failed to delete task %d in cache: %v", id, err)
	}

	return patched, nil
}
