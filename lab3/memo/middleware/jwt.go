package middleware

import (
	"context"
	"memo/pkg/util"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func JWT(ctx context.Context, r *app.RequestContext) {
	token := string(r.GetHeader("Authorization"))
	if token == "" {
		r.JSON(http.StatusBadRequest, 0) // TODO
		r.Abort()
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		// TODO: 解析token失败
		r.Abort()
	} else if time.Now().Unix() > claims.ExpiresAt.Unix() {
		// TODO: token过期
		r.Abort()
	} else {
		r.Set("uid", claims.UserId)
		r.Next(ctx)
	}
}
