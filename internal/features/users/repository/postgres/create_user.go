package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"practice/internal/core/domain"
	core_error "practice/internal/core/errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (id, full_name, phone_number)
	VALUES ($1, $2, $3)
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.ID, user.FullName, user.PhoneNumber)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.User{}, fmt.Errorf("user already exists: %w", core_error.ErrAlreadyExists)
		}
		return domain.User{}, fmt.Errorf("scan err: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)
	return userDomain, nil
}
