package statistics_transport_grpc

import (
	"context"
	"practice/internal/core/domain"
	v1 "practice/pkg/pb/v1"
	"time"
)

type StatisticsGRPC struct {
	service StatisticsService
	v1.UnimplementedStatisticServiceV1Server
}

func NewStatisticsGRPC(service StatisticsService) *StatisticsGRPC {
	return &StatisticsGRPC{service: service}
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from, to *time.Time,
	) (domain.Statistics, error)
}
