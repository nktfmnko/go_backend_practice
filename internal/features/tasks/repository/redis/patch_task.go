package tasks_redis_repository

import (
	"context"
	"practice/internal/core/domain"
)

func (r *CachedTasksRepository) PatchTask(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
	patched, err := r.next.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, err
	}

	_ = r.rdb.Del(ctx, r.getTaskKey(id)).Err()

	return patched, nil
}
