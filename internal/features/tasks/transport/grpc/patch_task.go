package tasks_transport_grpc

import (
	"context"
	"practice/internal/core/domain"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *TasksGRPC) PatchTask(ctx context.Context, req *v1.PatchTaskRequest) (*v1.PatchTaskResponse, error) {
	if req.GetPatch() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "patch cannot be nil")
	}

	patchTask := domain.NewTaskPatchFromProto(req)
	task, err := g.service.PatchTask(ctx, int(req.GetId()), patchTask)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to patch user")
	}

	return &v1.PatchTaskResponse{Task: task.ToProto()}, nil
}
