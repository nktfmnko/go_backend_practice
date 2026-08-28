package users_transport_http

import (
	"context"
	"net/http"
	"practice/internal/core/domain"
	core_http_server "practice/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
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
}

func NewUsersHTTPHandler(service UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{usersService: service}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
	}
}
