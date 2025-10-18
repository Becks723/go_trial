package main

import (
	"memo/controller"
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
	return controller.NewController(
		service.NewUserService(repository.NewUserRepo()),
		service.NewMemoService(repository.NewMemoRepo()))
}

func initializeRouter(c *controller.Controller) *gin.Engine {
	router := gin.Default()
	router.POST("/user/signup", c.UserSignup)
	router.POST("/user/login", c.UserLogin)
	return router
}
