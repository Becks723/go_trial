package util

import (
	"errors"
	"time"

	"StreamCore/biz/domain"
	"github.com/golang-jwt/jwt/v5"
)

type userCustomClaims struct {
	UserId   uint
	Username string
	jwt.RegisteredClaims
}

func GenerateAccessToken(u *domain.User, secret string, expiresIn time.Duration) (result string, err error) {
	exp := time.Now().Add(expiresIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userCustomClaims{
		UserId:   u.Id,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "StreamCore",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(u *domain.User, secret string, expiresIn time.Duration) (result string, err error) {
	exp := time.Now().Add(expiresIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userCustomClaims{
		UserId: u.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "StreamCore",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})
	return token.SignedString([]byte(secret))
}

func HoursOf(n int) time.Duration {
	return time.Hour * time.Duration(n)
}

func ParseToken(token string, secret string) (claims *userCustomClaims, err error) {
	tk, err := jwt.ParseWithClaims(token, &userCustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return
	}

	claims, ok := tk.Claims.(*userCustomClaims)
	if !ok {
		err = errors.New("unknown claims type")
		return
	}
	return claims, nil
}
