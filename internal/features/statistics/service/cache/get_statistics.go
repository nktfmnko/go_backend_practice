package statistics_cache_service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"practice/internal/core/domain"
	"time"
)

func (s StatisticsCacheService) GetStatistics(ctx context.Context, userID *int, from, to *time.Time) (domain.Statistics, error) {
	key := s.generateCacheKey(userID, from, to)
	val, err := s.rdb.Get(ctx, key).Result()
	if err == nil {
		var stat domain.Statistics
		if err := json.Unmarshal([]byte(val), &stat); err != nil {
			return domain.Statistics{}, fmt.Errorf("err in unmarshal stat %w", err)
		}
		return stat, nil
	}
	stats, err := s.next.GetStatistics(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, err
	}

	data, err := json.Marshal(stats)
	if err != nil {
		log.Printf("ERROR: failed to marshal stat: %v", err)
	}

	if err := s.rdb.Set(ctx, key, data, s.ttl).Err(); err != nil {
		log.Printf("ERROR: failed to save stat to redis: %v", err)
	}

	return stats, nil
}
