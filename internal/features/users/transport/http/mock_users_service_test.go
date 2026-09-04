package users_transport_http

import (
	"context"
	"practice/internal/core/domain"
)

type mockUsersService struct {
	getUserFunc    func(ctx context.Context, id int) (domain.User, error)
	deleteUserFunc func(ctx context.Context, id int) error
	createUserFunc func(ctx context.Context, fullName string, phoneNumber *string) (domain.User, error)
	getUsersFunc   func(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)
	patchUserFunc func(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

func (m *mockUsersService) GetUser(ctx context.Context, id int) (domain.User, error) {
	return m.getUserFunc(ctx, id)
}

func (m *mockUsersService) DeleteUser(ctx context.Context, id int) error {
	return m.deleteUserFunc(ctx, id)
}

func (m *mockUsersService) CreateUser(ctx context.Context, fullName string, phoneNumber *string) (domain.User, error) {
	return m.createUserFunc(ctx, fullName, phoneNumber)
}

func (m *mockUsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	return m.getUsersFunc(ctx, limit, offset)
}

func (m *mockUsersService) PatchUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	return m.patchUserFunc(ctx, id, patch)
}
