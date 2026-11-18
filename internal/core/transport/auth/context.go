package core_auth

import (
	"context"
	core_error "practice/internal/core/errors"
)

type ctxKey string

const userClaimsKey ctxKey = "user_claims"

func WithUser(ctx context.Context, u UserClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, u)
}

func UserFromContext(ctx context.Context) (UserClaims, bool) {
	u, ok := ctx.Value(userClaimsKey).(UserClaims)
	return u, ok
}

func MustUserFromContext(ctx context.Context) (UserClaims, error) {
	claims, ok := UserFromContext(ctx)
	if !ok {
		return UserClaims{}, core_error.ErrUnauthenticated
	}
	return claims, nil
}
