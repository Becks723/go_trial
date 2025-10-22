package main

import (
	"memo/controller"
	_ "memo/docs"
	"memo/middleware"
	"memo/repository"
	"memo/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func initializeRouter(c *controller.Controller) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := router.Group("/api/v1")
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
	return router
}
