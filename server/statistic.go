package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/BUGLAN/stress/model"
)

func ReceiveData(ch chan *model.ReqResult, wg *sync.WaitGroup) {
	defer wg.Done()
	stopChan := make(chan struct{})
	ticker := time.NewTicker(time.Second * 1)

	// 定义指标
	var (
		qps              float64 // qps 每秒请求数
		requestTotalTime uint64  // 请求总时间
		totalProcessTime uint64  // 请求总耗时
		totalSuccessNum  uint64  // 请求总成功数
		totalFailureNum  uint64  // 请求总失败数
		MaxTime          uint64  // 单个请求最大耗时
		MinTime          uint64  // 单个请求最少耗时
		avgTime          uint64  // 平均请求耗时
		concurrentNum    uint64  // 并发数
		currRequestNum   uint64  // 当前请求数
	)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Printf("%4ds│%7d│%7d│%7d│%8.2f│%8d│%8d│%8d│ 错误码\n", totalProcessTime/1e9, concurrentNum, totalSuccessNum, totalFailureNum, qps, MaxTime/1e9, MinTime/1e9, avgTime/1e9)
				// tableContent()
			case <-stopChan:
				fmt.Println("STOP!!!")
				return
			}
		}

	}()

	// 输出表头
	tableHeader()

	// 接收channel中的数据
	for data := range ch {
		// 最大耗时
		if data.ProcessTime > MaxTime {
			MaxTime = data.ProcessTime
		}

		// 最小耗时
		if data.ProcessTime < MinTime && MinTime != 0 {
			MinTime = data.ProcessTime
		}

		// 请求总耗时
		totalProcessTime += data.ProcessTime

		// 请求成功/失败数
		if data.IsSuccess {
			totalSuccessNum++
		} else {
			totalFailureNum++
		}

		// 并发数
		concurrentNum++

		// 平均耗时
		avgTime = totalProcessTime / concurrentNum

		// qps
		qps = float64(1e9 / avgTime)

	}

	// channel中数据已完成, 结束
	stopChan <- struct{}{}

	// 输出压测文档
	// done
	_ = requestTotalTime
	_ = totalSuccessNum
	_ = totalFailureNum
	_ = avgTime
	_ = concurrentNum
	_ = currRequestNum
	fmt.Printf("%4ds│%7d│%7d│%7d│%8.2f│%8d│%8d│%8d│错误码\n", totalProcessTime/1e9, concurrentNum, totalSuccessNum, totalFailureNum, qps, MaxTime/1e9, MinTime/1e9, avgTime)

}

func tableHeader() {
	fmt.Printf("\n")
	fmt.Println(" 耗时│ 并发数│ 成功数│ 失败数│   QPS  │最长耗时│最短耗时│平均耗时│ 错误码")
	fmt.Println("─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────")
}

func table() {

}

func work() {

}
