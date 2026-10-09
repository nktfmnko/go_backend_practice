package users_transport_grpc

import (
	"context"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"
)

func (g *UsersGRPC) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	user, err := g.usersService.GetUser(ctx, int(req.GetId()))
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to get user")
	}

	return &v1.GetUserResponse{User: user.ToProtoUser()}, nil
}
