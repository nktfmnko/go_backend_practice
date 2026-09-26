package statistics_cache_service

import (
	"fmt"
	core_redis "practice/internal/core/repository/cache/redis"
	statistics_transport_http "practice/internal/features/statistics/transport/http"
	"time"
)

type StatisticsCacheService struct {
	next statistics_transport_http.StatisticsService
	rdb  *core_redis.RedisClient
	ttl  time.Duration
}

func NewStatisticsCacheService(
	next statistics_transport_http.StatisticsService,
	rdb *core_redis.RedisClient,
	ttl time.Duration) *StatisticsCacheService {
	return &StatisticsCacheService{
		next: next,
		rdb:  rdb,
		ttl:  ttl,
	}
}

func (s StatisticsCacheService) generateCacheKey(userID *int, from, to *time.Time) string {
	userVal := "all"
	if userID != nil {
		userVal = fmt.Sprintf("%d", *userID)
	}

	fromVal := "min"
	if from != nil {
		fromVal = from.UTC().Format(time.RFC3339)
	}

	toVal := "max"
	if to != nil {
		toVal = to.UTC().Format(time.RFC3339)
	}

	return fmt.Sprintf("stats:user:%s:from:%s:to:%s", userVal, fromVal, toVal)
}
