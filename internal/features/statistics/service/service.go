package statistics_service

import (
	"context"
	"practice/internal/core/domain"
	"time"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

func NewStatisticsService(statisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{statisticsRepository: statisticsRepository}
}

type StatisticsRepository interface {
	GetTasks(ctx context.Context, userID *int, from, to *time.Time) ([]domain.Task, error)
}
