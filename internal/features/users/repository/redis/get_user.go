package users_redis_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"practice/internal/core/domain"

	"github.com/redis/go-redis/v9"
)

func (r *CachedUsersRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	key := r.getUserKey(id)
	val, err := r.rdb.Get(ctx, key).Result()
	if err == nil {
		var user domain.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			return user, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		fmt.Printf("Redis error: %v", err)
	}

	user, err := r.next.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	if data, err := json.Marshal(user); err == nil {
		_ = r.rdb.Set(ctx, key, data, r.ttl).Err()
	}

	return user, nil
}
