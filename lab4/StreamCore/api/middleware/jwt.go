package middleware

import (
	"StreamCore/api/pack"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util/jwt"
	"context"
	"errors"
	"time"

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

		env := env.Instance()
		claims, err := jwt.ParseToken(access, env.AccessToken_Secret)
		// access failed
		if err != nil ||
			time.Now().Unix() > claims.ExpiresAt.Unix() {
			// refresh access
			newAccess, newRefresh, err := refreshToken(refresh, env.RefreshToken_Secret)
			if err != nil {
				pack.RespUnauthorizedError(c, err)
				c.Abort()
				return
			}
			// update headers
			c.Header(AccessTokenKey, newAccess)
			c.Header(RefreshTokenKey, newRefresh)
			claims, _ = jwt.ParseToken(newAccess, env.AccessToken_Secret)
		}
		c.Set("uid", claims.UserId)
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

	env := env.Instance()
	// new access token
	newAccess, err = jwt.GenerateAccessToken(claims.UserId, env.AccessToken_Secret, jwt.HoursOf(env.AccessToken_ExpiryHours))
	if err != nil {
		return "", "", err
	}
	// new refresh token (optional)
	newRefresh, err = jwt.GenerateRefreshToken(claims.UserId, env.RefreshToken_Secret, jwt.HoursOf(env.RefreshToken_ExpiryHours))
	if err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}
