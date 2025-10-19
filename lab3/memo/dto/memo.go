package dto

import "time"

type AddMemoReq struct {
	Title    string     `form:"title"`
	Content  string     `form:"content"`
	StartsAt *time.Time `form:"starts_at"` // RFC3339: 2000-01-01T23:59:59Z
	EndsAt   *time.Time `form:"ends_at"`
}
