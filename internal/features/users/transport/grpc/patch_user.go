package users_transport_grpc

import (
	"context"
	"practice/internal/core/domain"
	core_auth "practice/internal/core/transport/auth"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *UsersGRPC) PatchUser(ctx context.Context, req *v1.PatchUserRequest) (*v1.PatchUserResponse, error) {
	claims, err := core_auth.MustUserFromContext(ctx)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "unauthenticated")
	}

	if req.GetPatch() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "patch cannot be nil")
	}
	patch := domain.NewUserPatchFromProto(req)

	user, err := g.usersService.PatchUser(ctx, claims.UserID, int(req.GetId()), patch)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to patch user")
	}

	return &v1.PatchUserResponse{User: user.ToProtoUser()}, nil
}
