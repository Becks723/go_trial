package controller

import (
	"memo/dto"
	"memo/pkg/ctl"
	"memo/pkg/e"
	"memo/service"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	userServ *service.UserService
	memoServ *service.MemoService
}

/* ctor for Controller struct. */
func NewController(userServ *service.UserService, memoServ *service.MemoService) *Controller {
	return &Controller{
		userServ: userServ,
		memoServ: memoServ,
	}
}

func (c *Controller) UserSignup(ctx *gin.Context) {
	var req dto.SignupReq
	if err := ctx.ShouldBind(&req); err == nil {
		resp, err := c.userServ.Signup(&req)
		if err == nil {
			ctx.JSON(e.Success, resp)
		} else {
			ctx.JSON(e.InternalError, ctl.ResponseError(err))
		}
		return
	} else {
		ctx.JSON(e.BadRequest, ctl.ResponseError(err, e.BadRequest))
	}
}

func (c *Controller) UserLogin(ctx *gin.Context) {
	var req dto.LoginReq
	if err := ctx.ShouldBind(&req); err == nil {
		resp, err := c.userServ.Login(&req)
		if err == nil {
			ctx.JSON(e.Success, resp)
		} else {
			ctx.JSON(e.InternalError, ctl.ResponseError(err))
		}
		return
	} else {
		ctx.JSON(e.BadRequest, ctl.ResponseError(err, e.BadRequest))
	}
}

func (c *Controller) MemoAdd(ctx *gin.Context) {

}

func (c *Controller) MemoUpdate(ctx *gin.Context) {

}

func (c *Controller) MemoFind(ctx *gin.Context) {

}

func (c *Controller) MemoDelete(ctx *gin.Context) {

}
