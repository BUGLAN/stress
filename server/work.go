package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/BUGLAN/stress/client"
	"github.com/BUGLAN/stress/model"
)

// Stress run the stress testing
func Stress(url string, ch chan *model.ReqResult, wg *sync.WaitGroup) {
	var err error
	for i := 0; i < model.Num; i++ {
		requestTime := time.Now().UnixNano()

		err = client.Get(url, http.Header{})

		isSuccess := true
		if err != nil {
			isSuccess = false
			if model.Debug {
				fmt.Printf("err: %s\n", err.Error())
			}
		}
		// 每个请求都开协程的话, 消耗过大, 反而不利于并发请求
		ch <- &model.ReqResult{
			IsSuccess:   isSuccess,
			StatusCode:  200,
			ProcessTime: float64(time.Now().UnixNano() - requestTime),
			RequestTime: uint64(requestTime),
		}
	}
	wg.Done()
}
