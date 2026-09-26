package users_redis_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"practice/internal/core/domain"
	"time"
)

func (r *CachedUsersRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	key := r.getUserKey(id)
	val, err := r.rdb.Get(ctx, key).Result()
	if err == nil {
		var user domain.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			return domain.User{}, fmt.Errorf("err in unmarshal user: %w", err)
		}
		return user, nil
	}
	res, err, _ := r.singleFlightGroup.Do(key, func() (interface{}, error) {
		opCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := r.next.GetUser(opCtx, id)
		if err != nil {
			r.singleFlightGroup.Forget(key)
			return domain.User{}, err
		}

		data, err := json.Marshal(user)
		if err != nil {
			log.Printf("ERROR: failed to marshal user %d for cache: %v", id, err)
			r.singleFlightGroup.Forget(key)
			return user, nil
		}

		if err := r.rdb.Set(opCtx, key, data, r.ttl).Err(); err != nil {
			log.Printf("ERROR: failed to save user %d to redis: %v", id, err)
		}
		return user, nil
	})

	if err != nil {
		return domain.User{}, err
	}
	user, ok := res.(domain.User)
	if !ok {
		return domain.User{}, fmt.Errorf("unexpected type in singleflight res: %T", res)
	}

	return user, nil
}
