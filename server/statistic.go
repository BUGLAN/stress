package server

import (
	"fmt"
	"math"
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
		totalProcessTime float64 // 请求总耗时
		totalSuccessNum  uint64  // 请求总成功数
		totalFailureNum  uint64  // 请求总失败数
		maxTime          float64 // 单个请求最大耗时
		minTime          float64 // 单个请求最少耗时
		avgTime          float64 // 平均请求耗时
		concurrentNum    uint64  // 并发数
		currRequestNum   uint64  // 当前请求数
	)

	// 排除minTime为0的情况
	minTime = math.MaxFloat64

	go func() {
		for {
			select {
			case <-ticker.C:
				// 定时输出相应的指标
				out(totalProcessTime, concurrentNum, totalSuccessNum, totalFailureNum, qps, maxTime, minTime, avgTime)
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
		if data.ProcessTime > maxTime {
			maxTime = data.ProcessTime
		}

		// fmt.Printf("data.ProcessTime %f\n", data.ProcessTime)
		// fmt.Printf("minTime %f\n", minTime)

		// 最小耗时
		if data.ProcessTime < minTime {
			minTime = data.ProcessTime
		}

		totalProcessTime += data.ProcessTime

		// 请求成功/失败数
		if data.IsSuccess {
			totalSuccessNum++
		} else {
			totalFailureNum++
		}

		// 并发数(请求数)
		concurrentNum++

		// 平均耗时
		avgTime = totalProcessTime / float64(concurrentNum)

		// qps
		qps = 1e9 / avgTime

	}

	// channel中数据已完成, 结束
	stopChan <- struct{}{}

	// 输出压测文档
	// done
	_ = requestTotalTime
	_ = concurrentNum
	_ = currRequestNum
	// 最后的输出到控制台 传入
	out(totalProcessTime, concurrentNum, totalSuccessNum, totalFailureNum, qps, maxTime, minTime, avgTime)

}

func tableHeader() {
	fmt.Printf("\n")
	fmt.Println(" 耗时│ 并发数│ 成功数│ 失败数│   QPS  │最长耗时│最短耗时│平均耗时│ 错误码")
	fmt.Println("─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────")
}

// out 输出到控制台 单位为纳秒(ns)
func out(totalProcessTime float64, concurrentNum, totalSuccessNum, totalFailureNum uint64, qps, maxTime, minTime, avgTime float64) {
	fmt.Printf("%4.2fs│%7d│%7d│%7d│%8.2f│%8.2fs│%8.2fs│%8.2fs│错误码\n", totalProcessTime/1e9, concurrentNum, totalSuccessNum, totalFailureNum, qps, maxTime/1e9, minTime/1e9, avgTime/1e9)

}
