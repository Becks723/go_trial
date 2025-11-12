package util

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
)

func RetrieveUserId(ctx context.Context) (uint, error) {
	obj := ctx.Value("uid")
	if obj != nil {
		uid, ok := obj.(uint)
		if ok {
			return uid, nil
		}
	}
	return 0, errors.New("Error retrieving uid.")
}

func ContextWithUid(ctx context.Context, c *app.RequestContext) context.Context {
	obj, _ := c.Get("uid")
	uid, _ := obj.(uint)
	return context.WithValue(ctx, "uid", uid)
}
