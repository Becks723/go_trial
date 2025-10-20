package main

import (
	"memo/controller"
	"memo/middleware"
	"memo/repository"
	"memo/service"

	"github.com/gin-gonic/gin"
)

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

func initializeRouter(c *controller.Controller) *gin.Engine {
	router := gin.Default()
	router.POST("/user/signup", c.UserSignup)
	router.POST("/user/login", c.UserLogin)

	authed := router.Group("/")
	authed.Use(middleware.JWT)
	{
		authed.POST("memo/add", c.MemoAdd)
		authed.POST("memo/update", c.MemoUpdate)
	}
	return router
}
