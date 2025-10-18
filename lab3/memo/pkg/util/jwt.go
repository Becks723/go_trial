package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = "my_jwt_secret"

// 自定义Claims
type UserClaims struct {
	UserId   uint
	Username string
	jwt.RegisteredClaims
}

// 生成token
func GenerateToken(id uint, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserId:   id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "memo",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})
	return token.SignedString([]byte(jwtSecret))
}

// 解析token信息
func ParseToken(token string) (claims *UserClaims, err error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return
	}
	claims, ok := tokenClaims.Claims.(*UserClaims)
	if ok {
		return
	} else {
		err = errors.New("Unknown claims type.") // TODO: i18n
		return
	}
}
