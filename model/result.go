package model

type ReqResult struct {
	StatusCode  int
	RequestTime uint64 // 请求时间
	ProcessTime float64 // 处理时间
	IsSuccess   bool   // 是否成功
}