package handler

import (
	"StreamCore/biz/router/user"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestRegister(t *testing.T) {
	h := server.Default()
	user.Register(h)
	ut.PerformRequest(h.Engine, "POST", "/user/register", &ut.Body{})
}
