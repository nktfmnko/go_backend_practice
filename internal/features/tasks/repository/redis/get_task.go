package tasks_redis_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"practice/internal/core/domain"

	"github.com/redis/go-redis/v9"
)

func (r *CachedTasksRepository) GetTask(ctx context.Context, id int) (domain.Task, error) {
	taskKey := r.getTaskKey(id)
	v, err := r.rdb.Get(ctx, taskKey).Result()
	if err == nil {
		var task domain.Task
		if err := json.Unmarshal([]byte(v), &task); err == nil {
			return task, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		fmt.Println("redis err")
	}

	task, err := r.next.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}

	if data, err := json.Marshal(task); err == nil {
		_ = r.rdb.Set(ctx, taskKey, data, r.ttl).Err()
	}

	return task, nil
}
