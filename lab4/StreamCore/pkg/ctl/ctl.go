package ctl

import (
	"StreamCore/biz/model/common"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func ResponseError(err error, code ...int) *common.RespStatus {
	c := consts.StatusInternalServerError
	if len(code) > 0 {
		c = code[0]
	}
	return &common.RespStatus{
		Code: int32(c),
		Msg:  err.Error(),
	}
}

func ResponseSuccess() *common.RespStatus {
	return &common.RespStatus{
		Code: consts.StatusOK,
		Msg:  "ok",
	}
}
