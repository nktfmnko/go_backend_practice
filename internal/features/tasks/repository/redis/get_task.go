package tasks_redis_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"practice/internal/core/domain"
	"time"
)

func (r *CachedTasksRepository) GetTask(ctx context.Context, id int) (domain.Task, error) {
	taskKey := r.getTaskKey(id)
	v, err := r.rdb.Get(ctx, taskKey).Result()
	if err == nil {
		var task domain.Task
		if err := json.Unmarshal([]byte(v), &task); err != nil {
			return domain.Task{}, fmt.Errorf("err unmarshal: %w", err)
		}
		return task, nil
	}

	res, err, _ := r.singleFlightGroup.Do(taskKey, func() (interface{}, error) {
		opCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		task, err := r.next.GetTask(opCtx, id)
		if err != nil {
			r.singleFlightGroup.Forget(taskKey)
			return domain.Task{}, err
		}

		data, err := json.Marshal(task)
		if err != nil {
			log.Printf("ERROR: failed to marshal task %d for cache: %v", id, err)
			r.singleFlightGroup.Forget(taskKey)
			return task, nil
		}
		if err := r.rdb.Set(ctx, taskKey, data, r.ttl).Err(); err != nil {
			log.Printf("ERROR: failed to save task %d to redis: %v", id, err)
		}

		return task, nil
	})

	if err != nil {
		return domain.Task{}, err
	}

	task, ok := res.(domain.Task)
	if !ok {
		return domain.Task{}, fmt.Errorf("unexpected type in res: %T", res)
	}

	return task, nil
}
