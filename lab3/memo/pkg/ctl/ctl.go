/*
pkg/ctl/ctl.go -- 处理controller层统一响应
*/
package ctl

import (
	"memo/dto"
	"memo/pkg/e"
	"net/http"
)

func ResponseSuccess(code ...int) *dto.Response {
	c := http.StatusOK
	if len(code) != 0 {
		c = code[0]
	}

	return &dto.Response{
		Code:    c,
		Message: "ok",
		Data:    nil,
	}
}

func ResponseSuccessWithData(data any, code ...int) *dto.Response {
	c := http.StatusOK
	if len(code) != 0 {
		c = code[0]
	}

	return &dto.Response{
		Code:    c,
		Message: "ok",
		Data:    data,
	}
}

func ResponseError(err error, code ...int) *dto.Response {
	c := e.InternalError
	if len(code) != 0 {
		c = code[0]
	}

	return &dto.Response{
		Code:    c,
		Message: "error",
		Data:    err.Error(),
	}
}
