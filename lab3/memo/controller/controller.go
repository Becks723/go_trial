package controller

import (
	"errors"
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
	var req dto.AddMemoReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(e.BadRequest, ctl.ResponseError(err, e.BadRequest))
		return
	}

	// retrieve current user id
	var uid uint
	var err error
	if uid, err = c.retrieveCurrentUid(ctx); err != nil {
		ctx.JSON(e.InternalError, ctl.ResponseError(err))
		return
	}

	// do add memo service
	var resp *dto.Response
	if resp, err = c.memoServ.Add(uid, &req); err != nil {
		ctx.JSON(e.InternalError, ctl.ResponseError(err))
		return
	}

	ctx.JSON(e.Success, resp)
}

func (c *Controller) MemoUpdate(ctx *gin.Context) {
	var req dto.UpdateMemoReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(e.BadRequest, ctl.ResponseError(err, e.BadRequest))
		return
	}

	// retrieve current user id
	var uid uint
	var err error
	if uid, err = c.retrieveCurrentUid(ctx); err != nil {
		ctx.JSON(e.InternalError, ctl.ResponseError(err))
		return
	}

	// do update memo service
	var resp *dto.Response
	if resp, err = c.memoServ.Update(uid, &req); err != nil {
		ctx.JSON(e.InternalError, ctl.ResponseError(err))
		return
	}

	ctx.JSON(e.Success, resp)
}

func (c *Controller) MemoList(ctx *gin.Context) {
	var params dto.ListMemoParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(e.BadRequest, ctl.ResponseError(err, e.BadRequest))
		return
	}

	// retrieve current user id
	var uid uint
	var err error
	if uid, err = c.retrieveCurrentUid(ctx); err != nil {
		ctx.JSON(e.InternalError, ctl.ResponseError(err))
		return
	}

	// do list memo service
	var resp *dto.Response
	if resp, err = c.memoServ.List(uid, &params); err != nil {
		ctx.JSON(e.InternalError, ctl.ResponseError(err))
		return
	}

	ctx.JSON(e.Success, resp)
}

func (c *Controller) MemoFind(ctx *gin.Context) {

}

func (c *Controller) MemoDelete(ctx *gin.Context) {

}

func (c *Controller) retrieveCurrentUid(ctx *gin.Context) (uid uint, err error) {
	raw, ok := ctx.Get("uid")
	if !ok {
		err = errors.New("Cannot retrieve current uid.") // TODO: i18n
	} else if uid, ok = raw.(uint); !ok {
		err = errors.New("Unknown type of uid.") // 不太可能出现，不翻译了
	}
	return
}
