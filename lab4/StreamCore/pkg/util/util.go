package util

import (
	"context"
	"errors"
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
