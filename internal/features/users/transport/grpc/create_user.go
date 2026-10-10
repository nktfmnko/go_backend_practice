package users_transport_grpc

import (
	"context"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *UsersGRPC) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	if req.GetInfo() == nil {
		return nil, status.Error(codes.InvalidArgument, "user info is required")
	}

	user, err := g.usersService.CreateUser(ctx, req.GetInfo().GetFullName(), req.GetInfo().PhoneNumber)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to create user")
	}

	return &v1.CreateUserResponse{User: user.ToProtoUser()}, nil
}
