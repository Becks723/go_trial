package middleware

import (
	"StreamCore/biz/domain"
	"StreamCore/biz/repo"
	"StreamCore/pkg/ctl"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util/jwt"
	"context"
	"errors"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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
			err := errors.New("Not authorized.")
			c.JSON(consts.StatusUnauthorized, ctl.ResponseError(err, consts.StatusUnauthorized))
			c.Abort()
			return
		}

		env := env.Instance()
		claims, err := jwt.ParseToken(access, env.AccessToken_Secret)
		// access failed
		if err != nil ||
			time.Now().Unix() > claims.ExpiresAt.Unix() {
			// refresh access
			newAccess, newRefresh, err := refreshToken(refresh, env.RefreshToken_Secret, repo.NewUserRepo())
			if err != nil {
				c.JSON(consts.StatusUnauthorized, ctl.ResponseError(err, consts.StatusUnauthorized))
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

func refreshToken(refresh string, secret string, ur repo.UserRepo) (newAccess, newRefresh string, err error) {
	// resolve refresh token
	claims, err := jwt.ParseToken(refresh, secret)
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
	newAccess, err = jwt.GenerateAccessToken(u, env.AccessToken_Secret, jwt.HoursOf(env.AccessToken_ExpiryHours))
	if err != nil {
		return "", "", err
	}
	// new refresh token (optional)
	newRefresh, err = jwt.GenerateRefreshToken(u, env.RefreshToken_Secret, jwt.HoursOf(env.RefreshToken_ExpiryHours))
	if err != nil {
		return "", "", err
	}
	return
}
