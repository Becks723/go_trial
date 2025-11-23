package service

import (
	"strconv"
	"time"
)

func parseTIme(timestamp string) (t time.Time, err error) {
	unix, err := strconv.ParseUint(timestamp, 10, 64)
	if err != nil {
		return
	}
	t = time.UnixMilli(int64(unix))
	return
}
