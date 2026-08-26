package users_service

import (
	"context"
	"practice/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

func NewUsersService(repository UsersRepository) *UsersService {
	return &UsersService{usersRepository: repository}
}
