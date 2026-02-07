package middleware

import (
	"context"
	"errors"
	"time"

	"StreamCore/api/pack"
	"StreamCore/internal/pkg/base/rpccontext"
	"StreamCore/internal/pkg/constants"
	"StreamCore/pkg/util/jwt"
	"github.com/cloudwego/hertz/pkg/app"
)

func JWTAuthFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		const (
			AccessTokenKey  = "Access-Token"
			RefreshTokenKey = "Refresh-Token"
		)
		access := string(c.GetHeader(AccessTokenKey))
		refresh := string(c.GetHeader(RefreshTokenKey))

		if access == "" && refresh == "" {
			pack.RespUnauthorizedError(c, errors.New("token not provided"))
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(access, constants.JWT_AccessSecret)
		// access failed
		if err != nil ||
			time.Now().Unix() > claims.ExpiresAt.Unix() {
			// refresh access
			newAccess, newRefresh, err := refreshToken(refresh, constants.JWT_RefreshSecret)
			if err != nil {
				pack.RespUnauthorizedError(c, err)
				c.Abort()
				return
			}
			// update headers
			c.Header(AccessTokenKey, newAccess)
			c.Header(RefreshTokenKey, newRefresh)
			claims, _ = jwt.ParseToken(newAccess, constants.JWT_AccessSecret)
		}
		ctx = rpccontext.WithLoginUid(ctx, claims.UserId)
		c.Next(ctx)
	}
}

func refreshToken(refresh string, secret string) (newAccess, newRefresh string, err error) {
	// resolve refresh token
	claims, err := jwt.ParseToken(refresh, secret)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// refresh expired
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return "", "", errors.New("refresh token expired")
	}

	// new access token
	newAccess, err = jwt.GenerateAccessToken(claims.UserId, constants.JWT_AccessSecret, constants.JWT_AccessTokenExpiration)
	if err != nil {
		return "", "", err
	}
	// new refresh token (optional)
	newRefresh, err = jwt.GenerateRefreshToken(claims.UserId, constants.JWT_RefreshSecret, constants.JWT_RefreshTokenExpiration)
	if err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}
