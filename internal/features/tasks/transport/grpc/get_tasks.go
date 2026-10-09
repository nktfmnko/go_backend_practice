package tasks_transport_grpc

import (
	"context"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"
)

func (g *TasksGRPC) GetTasks(ctx context.Context, req *v1.GetTasksRequest) (*v1.GetTasksResponse, error) {
	var userId, lim, off *int

	if userReq := req.GetUserId(); userReq != nil {
		val := int(userReq.GetValue())
		userId = &val
	}
	if limReq := req.GetLimit(); limReq != nil {
		val := int(limReq.GetValue())
		lim = &val
	}

	if offReq := req.GetLimit(); offReq != nil {
		val := int(offReq.GetValue())
		off = &val
	}

	tasks, err := g.service.GetTasks(ctx, userId, lim, off)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to get tasks")
	}

	resp := make([]*v1.Task, len(tasks))
	for i, t := range tasks {
		resp[i] = t.ToProto()
	}

	return &v1.GetTasksResponse{Tasks: resp}, nil
}
