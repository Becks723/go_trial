package dto

import (
	"memo/repository/model"
	"time"
)

type AddMemoReq struct {
	Title    string     `form:"title"`
	Content  string     `form:"content"`
	StartsAt *time.Time `form:"starts_at"` // RFC3339: 2000-01-01T23:59:59Z
	EndsAt   *time.Time `form:"ends_at"`
}

type UpdateMemoReq struct {
	Id       uint             `form:"id"`
	Title    string           `form:"title"`
	Content  string           `form:"content"`
	Status   model.MemoStatus `form:"status"`
	StartsAt *time.Time       `form:"starts_at"`
	EndsAt   *time.Time       `form:"ends_at"`
}
