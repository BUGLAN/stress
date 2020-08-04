package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/BUGLAN/stress/model"
)

func (srv *Server) Receiver(coroutineNum uint64, ch chan *model.ReqResult, wg *sync.WaitGroup) {
	defer wg.Done()
	stopChan := make(chan struct{})
	ticker := time.NewTicker(time.Second * 1)

	// 定义指标
	var (
		qps              float64 // qps 每秒请求数
		totalRequestTime float64 // 请求总耗时
		totalSuccessNum  uint64  // 请求总成功数
		totalFailureNum  uint64  // 请求总失败数
		maxTime          float64 // 单个请求最大耗时
		minTime          float64 // 单个请求最少耗时
		avgTime          float64 // 平均请求耗时
		concurrentNum    uint64  // 并发数
		currTime         float64 // 当前时间
	)

	// 定义startTime
	startTime := time.Now().UnixNano()

	go func() {
		for {
			select {
			case <-ticker.C:
				currTime = float64(time.Now().UnixNano() - startTime)
				out(currTime, concurrentNum, totalSuccessNum, totalFailureNum, qps, maxTime, minTime, avgTime)
			case <-stopChan:
				fmt.Println()
				return
			}
		}
	}()

	// 输出表头
	tableHeader()

	// 接收channel中的数据
	for data := range ch {
		// 最大耗时
		if data.RequestTime > maxTime {
			maxTime = data.RequestTime
		}

		// 最小耗时
		if minTime == 0 {
			minTime = data.RequestTime
		} else if data.RequestTime < minTime {
			minTime = data.RequestTime
		}

		totalRequestTime += data.RequestTime

		// 请求成功/失败数
		if data.IsSuccess {
			totalSuccessNum++
		} else {
			totalFailureNum++
		}

		// 请求数
		concurrentNum++

		// 平均耗时
		avgTime = totalRequestTime / float64(concurrentNum)

		// qps
		qps = float64(totalSuccessNum*1e9*coroutineNum) / totalRequestTime
	}

	// channel中数据已完成, 结束
	stopChan <- struct{}{}

	// 最后的输出到控制台 传入
	currTime = float64(time.Now().UnixNano() - startTime)
	out(currTime, concurrentNum, totalSuccessNum, totalFailureNum, qps, maxTime, minTime, avgTime)
}

func tableHeader() {
	fmt.Printf("\n")
	fmt.Println(" 耗时│ 请求数│ 成功数│ 失败数│   QPS  │ 最长耗时│ 最短耗时│ 平均耗时")
	fmt.Println("─────┼───────┼───────┼───────┼────────┼─────────┼─────────┼─────────")
}

// out 输出到控制台 单位为纳秒(ns)
func out(requestTime float64, concurrentNum, totalSuccessNum, totalFailureNum uint64, qps, maxTime, minTime, avgTime float64) {
	fmt.Printf("%4.0fs│%7d│%7d│%7d│%8.2f│%7.2fms│%7.2fms│%7.2fms \n", requestTime/1e9, concurrentNum, totalSuccessNum, totalFailureNum, qps, maxTime/1e6, minTime/1e6, avgTime/1e6)
}
