package core_auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrWrongApp     = errors.New("token issued for another app")
)

type UserClaims struct {
	UserID int64
	Email  string
}

func ParseToken(tokenString, appSecret string, expectedAppID int) (UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appSecret), nil
	})

	if err != nil || !token.Valid {
		return UserClaims{}, fmt.Errorf("%w:%v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return UserClaims{}, ErrInvalidToken
	}

	appIDFloat, ok := claims["app_id"].(float64)
	if !ok || int(appIDFloat) != expectedAppID {
		return UserClaims{}, ErrWrongApp
	}

	uidFloat, ok := claims["uid"].(float64)
	if !ok {
		return UserClaims{}, ErrInvalidToken
	}

	email, _ := claims["email"].(string)

	return UserClaims{
		UserID: int64(uidFloat),
		Email:  email,
	}, nil
}
