package dto

import "time"

type AddMemoReq struct {
	Title    string
	Content  string
	StartsAt *time.Time
	EndsAt   *time.Time
}
