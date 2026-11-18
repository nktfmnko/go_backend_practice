package users_transport_grpc

import (
	"context"
	core_auth "practice/internal/core/transport/auth"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *UsersGRPC) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	claims, err := core_auth.MustUserFromContext(ctx)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "unauthenticated")
	}

	if err := g.usersService.DeleteUser(ctx, claims.UserID, int(req.GetId())); err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to delete user")
	}
	return &emptypb.Empty{}, nil
}
