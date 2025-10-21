package dto

import (
	"memo/repository/model"
	"time"
)

type MemoData struct {
	Id       uint             `json:"id"`
	Title    string           `json:"title"`
	Content  string           `json:"content"`
	Status   model.MemoStatus `json:"status"`
	StartsAt *time.Time       `json:"starts_at"`
	EndsAt   *time.Time       `json:"ends_at"`
}

type AddMemoReq struct {
	Title    string     `form:"title" binding:"required"`
	Content  string     `form:"content" binding:"required"`
	StartsAt *time.Time `form:"starts_at" binding:"required"` // RFC3339: 2000-01-01T23:59:59Z
	EndsAt   *time.Time `form:"ends_at" binding:"required"`
}

type UpdateMemoReq struct {
	Id       uint             `form:"id" binding:"required"`
	Title    string           `form:"title"`
	Content  string           `form:"content"`
	Status   model.MemoStatus `form:"status"`
	StartsAt *time.Time       `form:"starts_at"`
	EndsAt   *time.Time       `form:"ends_at"`
}

type ListMemoParams struct {
	Limit     int        `form:"limit"` // 每页容量，url请求参数也依然用form标签，因为和表单一样，本质上都是 key1=val1&key2=val2 的格式
	PageStart int        `form:"ps"`    // 从第几页开始
	PageEnd   int        `form:"pe"`    // 到第几页结束
	Filter    ListFilter `form:"filter"`
}

type ListFilter int

const (
	ListFilterNone      ListFilter = iota // 所有
	ListFilterCompleted                   // 已完成
	ListFilterPending                     // 代办
)

type SearchMemoParams struct {
	Limit     int    `form:"limit"`                       // 每页容量
	PageStart int    `form:"ps"`                          // 从第几页开始
	PageEnd   int    `form:"pe"`                          // 到第几页结束
	Keywords  string `form:"keywords" binding:"required"` // 关键词
}

type DeleteMemoByIdReq struct {
	Id uint `form:"id" binding:"required"`
}

type DeleteMemoByFilterReq struct {
	Filter DeleteFilter `form:"filter" binding:"required"`
}
type DeleteFilter int

const (
	DeleteFilterNone      DeleteFilter = iota // 所有
	DeleteFilterCompleted                     // 已完成
	DeleteFilterPending                       // 代办
)
