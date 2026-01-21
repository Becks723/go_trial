package logincontext

import (
	"context"
	"errors"
	"fmt"
)

func RetrieveLoginUid(ctx context.Context) (uint, error) {
	obj := ctx.Value("uid")
	if obj == nil {
		return 0, errors.New("context.Value(\"uid\") failed")
	}
	uid, ok := obj.(uint)
	if !ok {
		return 0, fmt.Errorf("failed conv %t to uint", obj)
	}
	return uid, nil
}

func WithLoginUid(ctx context.Context, uid uint) context.Context {
	return context.WithValue(ctx, "uid", uid)
}
