package middleware

import (
	"StreamCore/biz/repo"
	"StreamCore/pkg/ctl"
	"StreamCore/pkg/env"
	"StreamCore/pkg/util"
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
		claims, err := util.ParseToken(access, env.AccessToken_Secret)
		// access failed
		if err != nil ||
			time.Now().Unix() > claims.ExpiresAt.Unix() {
			// refresh access
			newAccess, newRefresh, err := util.RefreshToken(refresh, env.RefreshToken_Secret, repo.NewUserRepo())
			if err != nil {
				c.JSON(consts.StatusUnauthorized, ctl.ResponseError(err, consts.StatusUnauthorized))
				c.Abort()
				return
			}
			// update headers
			c.Header(AccessTokenKey, newAccess)
			c.Header(RefreshTokenKey, newRefresh)
			claims, _ = util.ParseToken(newAccess, env.AccessToken_Secret)
		}
		c.Set("uid", claims.UserId)
		c.Next(ctx)
	}
}
