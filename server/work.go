package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/BUGLAN/stress/client"
	"github.com/BUGLAN/stress/model"
)

func Stress(url string, ch chan *model.ReqResult, wg *sync.WaitGroup) {
	defer wg.Done()
	var w sync.WaitGroup
	for i := 0; i < model.Num; i++ {
		w.Add(1)
		requestTime := time.Now().UnixNano()
		go Process(url, ch, requestTime, &w)
	}
	w.Wait()
}

func Process(url string, ch chan *model.ReqResult, requestTime int64, wg *sync.WaitGroup) {
	defer wg.Done()
	err := client.Get(url, http.Header{})

	isSuccess := true
	if err != nil {
		isSuccess = false
		if model.Debug {
			fmt.Printf("err: %s\n", err.Error())
		}
	}

	ch <- &model.ReqResult{
		IsSuccess:   isSuccess,
		StatusCode:  200,
		ProcessTime: uint64(time.Now().UnixNano() - requestTime),
		RequestTime: uint64(requestTime),
	}
}
