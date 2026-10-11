package statistics_transport_grpc

import (
	"context"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"
)

func (g *StatisticsGRPC) GetStatistics(ctx context.Context, req *v1.GetStatisticsRequest) (*v1.GetStatisticsResponse, error) {
	stat, err := g.service.GetStatistics(ctx,
		core_transport_grpc.PtrConvert[int64, int](req.UserId),
		core_transport_grpc.ConvertProtoToTimePtr(req.From),
		core_transport_grpc.ConvertProtoToTimePtr(req.To),
	)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to get statistics")
	}

	return &v1.GetStatisticsResponse{
		TasksCreated:               int64(stat.TasksCreated),
		TasksCompleted:             int64(stat.TasksCompleted),
		TaskCompletedRate:          core_transport_grpc.PtrConvert[float64, float32](stat.TasksCompletedRate),
		TasksAverageCompletionTime: core_transport_grpc.ConvertDurationToProto(stat.TasksAverageCompletionTime),
	}, nil
}
