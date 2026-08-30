package tasks_transport_http

import (
	"context"
	"net/http"
	"practice/internal/core/domain"
	core_http_server "practice/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

func NewTasksHTTPHandler(tasksService TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{tasksService: tasksService}
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
	}
}
