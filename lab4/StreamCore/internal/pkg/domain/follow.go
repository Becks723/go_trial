package domain

import "time"

type Follow struct {
	TargetUid uint
	StartedAt time.Time
}
