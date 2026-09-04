package users_service

import (
	"context"
	"practice/internal/core/domain"
)

type mockUsersRepository struct {
	getUserFunc    func(ctx context.Context, id int) (domain.User, error)
	deleteUserFunc func(ctx context.Context, id int) error
	createUserFunc func(ctx context.Context, user domain.User) (domain.User, error)
	getUsersFunc   func(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)
	patchUserFunc func(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)
}

func (r *mockUsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	return r.createUserFunc(ctx, user)
}

func (r *mockUsersRepository) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	return r.getUsersFunc(ctx, limit, offset)
}

func (r *mockUsersRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	return r.getUserFunc(ctx, id)
}

func (r *mockUsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	return r.deleteUserFunc(ctx, id)
}

func (r *mockUsersRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	return r.patchUserFunc(ctx, id, user)
}
