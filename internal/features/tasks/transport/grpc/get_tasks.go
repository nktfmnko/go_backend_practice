package tasks_transport_grpc

import (
	"context"
	"practice/internal/core/domain"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"
)

func (g *TasksGRPC) GetTasks(ctx context.Context, req *v1.GetTasksRequest) (*v1.GetTasksResponse, error) {
	tasks, err := g.service.GetTasks(ctx,
		core_transport_grpc.PtrInt64ToInt(req.UserId),
		core_transport_grpc.PtrInt64ToInt(req.Limit),
		core_transport_grpc.PtrInt64ToInt(req.Offset),
	)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to get tasks")
	}

	resp := domain.NewProtoSliceFromTasks(tasks)

	return &v1.GetTasksResponse{Tasks: resp}, nil
}
