package users_transport_grpc

import (
	"context"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"
)

func (g *UsersGRPC) GetUsers(ctx context.Context, req *v1.GetUsersRequest) (*v1.GetUsersResponse, error) {
	var lim, off *int
	if req.GetLimit() != nil {
		val := int(req.GetLimit().GetValue())
		lim = &val
	}
	if req.GetOffset() != nil {
		val := int(req.GetOffset().GetValue())
		off = &val
	}

	users, err := g.usersService.GetUsers(ctx, lim, off)
	if err != nil {
		return nil, core_transport_grpc.ErrorResponse(err, "failed to get users")
	}

	resp := make([]*v1.User, len(users))
	for i, u := range users {
		resp[i] = u.ToProtoUser()
	}

	return &v1.GetUsersResponse{Users: resp}, nil
}
