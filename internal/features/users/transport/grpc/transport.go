package users_transport_grpc

import (
	"context"
	"practice/internal/core/domain"
	v1 "practice/pkg/pb/v1"
)

type UsersGRPC struct {
	usersService UsersService
	v1.UnimplementedUsersServiceV1Server
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		fullName string,
		phoneNumber *string,
	) (domain.User, error)
	GetUsers(
		ctx context.Context,
		limit, offset *int,
	) ([]domain.User, error)
	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)
	DeleteUser(
		ctx context.Context,
		id int,
	) error
	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

func NewUsersGRPC(service UsersService) *UsersGRPC {
	return &UsersGRPC{
		usersService: service,
	}
}
