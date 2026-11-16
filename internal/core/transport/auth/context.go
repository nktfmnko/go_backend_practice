package core_auth

import "context"

type ctxKey string

const userClaimsKey ctxKey = "user_claims"

func WithUser(ctx context.Context, u UserClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, u)
}

func UserFromContext(ctx context.Context) (UserClaims, bool) {
	u, ok := ctx.Value(userClaimsKey).(UserClaims)
	return u, ok
}
