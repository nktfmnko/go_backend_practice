package tasks_transport_grpc

import (
	"context"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *TasksGRPC) CreateTask(ctx context.Context, req *v1.CreateTaskRequest) (*v1.CreateTaskResponse, error) {
	if req.GetInfo() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "info cannot be null")
	}

	task, err := g.service.CreateTask(ctx, req.Info.GetTitle(), req.Info.Description)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to create task")
	}

	return &v1.CreateTaskResponse{Task: task.ToProto()}, nil
}
