package tasks_transport_http

import (
	"net/http"
	core_logger "practice/internal/core/logger"
	core_http_request "practice/internal/core/transport/http/request"
	core_http_response "practice/internal/core/transport/http/response"
)

func (h *TasksHTTPHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid id")
		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(err, "fail to delete task")
		return
	}

	responseHandler.NoContentResponse()
}
