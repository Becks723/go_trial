package util

import (
	"strconv"
	"time"
)

func TimePtr2String(t *time.Time) string {
	if t != nil {
		return t.String()
	}
	return ""
}

func Uint2String(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}
