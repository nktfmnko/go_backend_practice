package users_redis_repository

import (
	"context"
	"fmt"
	"practice/internal/core/domain"
	core_redis "practice/internal/core/repository/cache/redis"
	users_service "practice/internal/features/users/service"
	"time"
)

type CachedUsersRepository struct {
	next users_service.UsersRepository
	rdb  *core_redis.RedisClient
	ttl  time.Duration
}

func (r *CachedUsersRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return r.next.CreateUser(ctx, user)
}

func (r *CachedUsersRepository) GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error) {
	return r.next.GetUsers(ctx, limit, offset)
}

func NewCachedUsersRepository(next users_service.UsersRepository, rdb *core_redis.RedisClient, ttl time.Duration) *CachedUsersRepository {
	return &CachedUsersRepository{
		next: next,
		rdb:  rdb,
		ttl:  ttl,
	}
}

func (r *CachedUsersRepository) getUserKey(id int) string {
	return fmt.Sprintf("user:%d", id)
}
