package tasks_redis_repository

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
	core_redis "practice/internal/core/repository/cache/redis"
	tasks_service "practice/internal/features/tasks/service"
	"time"

	"golang.org/x/sync/singleflight"
)

type CachedTasksRepository struct {
	next              tasks_service.TasksRepository
	rdb               *core_redis.RedisClient
	singleFlightGroup singleflight.Group
	ttl               time.Duration
}

func (r *CachedTasksRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	return r.next.CreateTask(ctx, task)
}

func (r *CachedTasksRepository) GetTasks(ctx context.Context, userID, limit, offset *int) ([]domain.Task, error) {
	return r.next.GetTasks(ctx, userID, limit, offset)
}

func NewCachedTasksRepository(next tasks_service.TasksRepository, rdb *core_redis.RedisClient, ttl time.Duration) *CachedTasksRepository {
	return &CachedTasksRepository{
		next: next,
		rdb:  rdb,
		ttl:  ttl,
	}
}

func (r *CachedTasksRepository) getTaskKey(id int) string {
	return fmt.Sprintf("task:%d", id)
}
