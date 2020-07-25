package model

type ReqResult struct {
	StatusCode  int
	RequestTime float64 // 处理时间
	IsSuccess   bool    // 是否成功
}
