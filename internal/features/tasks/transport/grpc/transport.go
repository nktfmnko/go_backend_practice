package tasks_transport_grpc

import (
	"context"
	"practice/internal/core/domain"
	v1 "practice/pkg/pb/v1"
)

type TasksGRPC struct {
	service TasksService
	v1.UnimplementedTasksServiceV1Server
}

func NewTasksGRPC(service TasksService) *TasksGRPC {
	return &TasksGRPC{service: service}
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
	GetTask(
		ctx context.Context,
		id int,
	) (domain.Task, error)
	GetTasks(
		ctx context.Context,
		userID, limit, offset *int,
	) ([]domain.Task, error)
	DeleteTask(
		ctx context.Context,
		id int,
	) error
	PatchTask(
		ctx context.Context,
		id int,
		patch domain.TaskPatch,
	) (domain.Task, error)
}
