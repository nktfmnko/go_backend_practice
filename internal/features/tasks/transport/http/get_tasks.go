package tasks_transport_http

import (
	"fmt"
	"net/http"
	core_logger "practice/internal/core/logger"
	core_http_request "practice/internal/core/transport/http/request"
	core_http_response "practice/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHTTPHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "get query params")
		return
	}

	tasks, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "get tasks")
		return
	}

	response := GetTasksResponse(tasksDTOFromDomains(tasks))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userKey   = "user_id"
		limitKey  = "limit"
		offsetKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get userID query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get limit query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get offset query param: %w", err)
	}
	return userID, limit, offset, nil
}
