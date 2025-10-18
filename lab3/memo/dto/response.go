package dto

/* 通用的响应结构体 */
type Response struct {
	Code    int    `json:"code"` // 状态码
	Message string `json:"msg"`  // 返回信息
	Data    any    `json:"data"` // 业务数据
}
