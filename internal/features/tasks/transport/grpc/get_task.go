package tasks_transport_grpc

import (
	"context"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"
)

func (g *TasksGRPC) GetTask(ctx context.Context, req *v1.GetTaskRequest) (*v1.GetTaskResponse, error) {
	task, err := g.service.GetTask(ctx, int(req.GetId()))
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to get task")
	}

	return &v1.GetTaskResponse{Task: task.ToProto()}, nil
}
