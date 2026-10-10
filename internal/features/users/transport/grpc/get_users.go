package users_transport_grpc

import (
	"context"
	"practice/internal/core/domain"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"
)

func (g *UsersGRPC) GetUsers(ctx context.Context, req *v1.GetUsersRequest) (*v1.GetUsersResponse, error) {
	users, err := g.usersService.GetUsers(ctx, core_transport_grpc.PtrInt64ToInt(req.Limit), core_transport_grpc.PtrInt64ToInt(req.Offset))
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to get users")
	}

	resp := domain.ProtoSliceFromUsers(users)

	return &v1.GetUsersResponse{Users: resp}, nil

}
