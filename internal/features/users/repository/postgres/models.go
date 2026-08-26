package users_postgres_repository

import "practice/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromModels(users []UserModel) []domain.User {
	usersDomain := make([]domain.User, len(users))
	for i, user := range users {
		usersDomain[i] = domain.NewUser(user.ID, user.Version, user.FullName, user.PhoneNumber)
	}
	return usersDomain
}
