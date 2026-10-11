package tasks_transport_grpc

import (
	"context"
	"practice/internal/core/domain"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *TasksGRPC) CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.CreateTaskResponse, error) {
	if req.GetInfo() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "info cannot be null")
	}

	domainTask := domain.NewTaskUninitialized(req.Info.GetTitle(), req.Info.Description, int(req.Info.GetAuthorUserId()))

	task, err := g.service.CreateTask(ctx, domainTask)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to create task")
	}

	return &v1.CreateTaskResponse{Task: task.ToProto()}, nil
}
