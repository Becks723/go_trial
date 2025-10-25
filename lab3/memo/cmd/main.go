package main

import (
	"memo/controller"
	_ "memo/docs"
	"memo/middleware"
	"memo/repository"
	"memo/service"

	"github.com/cloudwego/hertz/pkg/app/server"
	hertzSwagger "github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title Memo API
// @version 1.0
// @description This is a memo api trial.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	repository.Load()

	ctl := newController()

	_ = initializeRouter(ctl).Run()
}

func newController() *controller.Controller {
	userRepo := repository.NewUserRepo()
	memoRepo := repository.NewMemoRepo()

	return controller.NewController(
		service.NewUserService(userRepo),
		service.NewMemoService(memoRepo, userRepo))
}

func initializeRouter(c *controller.Controller) *server.Hertz {
	h := server.Default(
		server.WithHostPorts("127.0.0.1:8080"))
	h.GET("/swagger/*any", hertzSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := h.Group("/api/v1")
	{
		v1.POST("/user/signup", c.UserSignup)
		v1.POST("/user/login", c.UserLogin)

		authed := v1.Group("/")
		authed.Use(middleware.JWT)
		{
			authed.POST("/memo/add", c.MemoAdd)
			authed.POST("/memo/update", c.MemoUpdate)
			authed.GET("/memo/list", c.MemoList)
			authed.GET("/memo/search", c.MemoSearch)
			authed.POST("/memo/delete", c.MemoDelete)
		}
	}
	return h
}
