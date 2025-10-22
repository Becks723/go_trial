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

// UserSignup godoc
// @Summary 用户注册
// @Description 用户注册
// @Tags user
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /user/signup [post]
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

// UserLogin godoc
// @Summary 用户登录
// @Description 用户登录
// @Tags user
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /user/login [post]
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

// MemoAdd godoc
// @Summary 创建代办事项
// @Description 创建代办事项
// @Tags memo
// @Accept x-www-form-urlencoded
// @Produce json
// @Param title formData string true "标题"
// @Param content formData string true "正文"
// @Param starts_at formData string true "开始时间，格式2000-01-01T23:59:59Z"
// @Param ends_at formData string true "截止时间，格式2000-01-01T23:59:59Z"
// @Header 200 {string} Authorization "必需，登录校验字段"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /memo/add [post]
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

// MemoUpdate godoc
// @Summary 更新代办事项
// @Description 更新代办事项
// @Tags memo
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id formData uint true "代办事项id"
// @Param title formData string false "标题"
// @Param content formData string false "正文"
// @Param starts_at formData *time.Time false "开始时间，格式2000-01-01T23:59:59Z"
// @Param ends_at formData *time.Time false "截止时间，格式2000-01-01T23:59:59Z"
// @Param status formData model.MemoStatus false "状态"
// @Header 200 {string} Authorization "必需，登录校验字段"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /memo/update [post]
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

// MemoList godoc
// @Summary 查询代办事项
// @Description 查询代办事项
// @Tags memo
// @Produce json
// @Param limit query int false "分页容量"
// @Param ps query int false "起始页码"
// @Param pe query int false "终止页码"
// @Param filter query dto.ListFilter false "筛选条件"
// @Header 200 {string} Authorization "必需，登录校验字段"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /memo/list [get]
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

// MemoSearch godoc
// @Summary 搜索代办事项
// @Description 搜索代办事项
// @Tags memo
// @Produce json
// @Param limit query int false "分页容量"
// @Param ps query int false "起始页码"
// @Param pe query int false "终止页码"
// @Param keywords query string true "必需，搜索关键词"
// @Header 200 {string} Authorization "必需，登录校验字段"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /memo/search [get]
func (c *Controller) MemoSearch(ctx *gin.Context) {
	var params dto.SearchMemoParams
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

	// do search memo service
	var resp *dto.Response
	if resp, err = c.memoServ.Search(uid, &params); err != nil {
		ctx.JSON(e.InternalError, ctl.ResponseError(err))
		return
	}

	ctx.JSON(e.Success, resp)
}

// MemoDelete godoc
// @Summary 删除代办事项
// @Description 删除代办事项。可删除指定id的事项，或者批量删除指定筛选条件的事项。
// @Tags memo
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id formData uint true "代办事项id"
// @Param filter formData dto.DeleteFilter true "筛选条件"
// @Header 200 {string} Authorization "必需，登录校验字段"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /memo/delete [post]
func (c *Controller) MemoDelete(ctx *gin.Context) {
	// retrieve current user id
	var uid uint
	var err error
	if uid, err = c.retrieveCurrentUid(ctx); err != nil {
		ctx.JSON(e.InternalError, ctl.ResponseError(err))
		return
	}

	var resp *dto.Response

	// try delete by id
	var ireq dto.DeleteMemoByIdReq
	if err = ctx.ShouldBind(&ireq); err == nil {
		// do delete by id service
		if resp, err = c.memoServ.DeleteById(uid, &ireq); err != nil {
			ctx.JSON(e.InternalError, ctl.ResponseError(err))
			return
		}
		ctx.JSON(e.Success, resp)
		return
	}

	// then try delete by filter
	var freq dto.DeleteMemoByFilterReq
	if err = ctx.ShouldBind(&freq); err == nil {
		// do delete by filter service
		if resp, err = c.memoServ.DeleteByFilter(uid, &freq); err != nil {
			ctx.JSON(e.InternalError, ctl.ResponseError(err))
			return
		}
		ctx.JSON(e.Success, resp)
		return
	}

	ctx.JSON(e.BadRequest, ctl.ResponseError(err, e.BadRequest))
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
