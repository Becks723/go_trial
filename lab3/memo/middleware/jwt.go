package middleware

import (
	"memo/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, 0) // TODO
		ctx.Abort()
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		// TODO: 解析token失败
		ctx.Abort()
	} else if time.Now().Unix() > claims.ExpiresAt.Unix() {
		// TODO: token过期
		ctx.Abort()
	} else {
		ctx.Set("uid", claims.UserId)
		ctx.Next()
	}
}
