package statistics_cache_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"practice/internal/core/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s StatisticsCacheService) GetStatistics(ctx context.Context, userID *int, from, to *time.Time) (domain.Statistics, error) {
	key := s.generateCacheKey(userID, from, to)

	val, err := s.rdb.Get(ctx, key).Result()
	if err == nil {
		var stat domain.Statistics
		if err := json.Unmarshal([]byte(val), &stat); err == nil {
			return stat, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		fmt.Println("Error in redis")
	}

	stats, err := s.next.GetStatistics(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, err
	}

	if data, err := json.Marshal(stats); err == nil {
		_ = s.rdb.Set(ctx, key, data, s.ttl).Err()
	}

	return stats, nil
}
