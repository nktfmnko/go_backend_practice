package tasks_transport_grpc

import (
	"context"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *TasksGRPC) DeleteTask(ctx context.Context, req *v1.DeleteTaskRequest) (*emptypb.Empty, error) {
	if err := g.service.DeleteTask(ctx, int(req.GetId())); err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to delete task")
	}

	return &emptypb.Empty{}, nil
}
