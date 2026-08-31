package statistics_transport_http

import (
	"context"
	"net/http"
	"practice/internal/core/domain"
	core_http_server "practice/internal/core/transport/http/server"
	"time"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

func NewStatisticsHTTPHandler(statisticsService StatisticsService) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{statisticsService: statisticsService}
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from, to *time.Time,
	) (domain.Statistics, error)
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
