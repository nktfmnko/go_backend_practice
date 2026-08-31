package statistics_transport_http

import (
	"fmt"
	"net/http"
	"practice/internal/core/domain"
	core_logger "practice/internal/core/logger"
	core_http_request "practice/internal/core/transport/http/request"
	core_http_response "practice/internal/core/transport/http/response"
	"time"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasksCreated"`
	TasksCompleted             int      `json:"tasksCompleted"`
	TasksCompletedRate         *float64 `json:"tasksCompletedRate"`
	TasksAverageCompletionTime *string  `json:"tasksAverageCompletionTime"`
}

func statisticsDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		avgTime = new(statistics.TasksAverageCompletionTime.String())
	}
	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func (h *StatisticsHTTPHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	userID, from, to, err := GetUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to get query params")
		return
	}

	stats, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to get statistics")
		return
	}

	response := statisticsDTOFromDomain(stats)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func GetUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDKey = "user_id"
		fromKey   = "from"
		toKey     = "to"
	)
	userID, err := core_http_request.GetIntQueryParam(r, userIDKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get userId query param:%w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get from query param:%w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get to query param:%w", err)
	}

	return userID, from, to, err
}
