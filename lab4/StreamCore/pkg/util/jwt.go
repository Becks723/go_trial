package util

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/repo"
	"StreamCore/pkg/env"
	"errors"
	"time"

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
		err = errors.New("Unknown claims type.")
		return
	}
	return claims, nil
}

func RefreshToken(refresh string, secret string, ur repo.UserRepo) (newAccess, newRefresh string, err error) {
	// resolve refresh token
	claims, err := ParseToken(refresh, secret)
	if err != nil {
		return "", "", errors.New("Error resolving refresh token.")
	}

	// refresh expired
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return "", "", errors.New("Refresh token expired.")
	}

	// check if user still exists
	var u *domain.User
	if u, err = ur.GetById(claims.UserId); err != nil {
		return "", "", errors.New("Refresh token: user not found.")
	}

	env := env.Instance()
	// new access token
	newAccess, err = GenerateAccessToken(u, env.AccessToken_Secret, HoursOf(env.AccessToken_ExpiryHours))
	if err != nil {
		return "", "", err
	}
	// new refresh token (optional)
	newRefresh, err = GenerateRefreshToken(u, env.RefreshToken_Secret, HoursOf(env.RefreshToken_ExpiryHours))
	if err != nil {
		return "", "", err
	}
	return
}
