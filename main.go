package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/BUGLAN/stress/client"
	"github.com/BUGLAN/stress/model"
)

var (
	coroutines int
	num        int
	help       bool
	url        string
	debug      bool
)

func main() {
	startTime := time.Now()

	// set runtime process
	runtime.GOMAXPROCS(1)

	// set flag vars
	flag.StringVar(&url, "u", "", "url 链接 https://www.baidu.com")
	flag.BoolVar(&help, "h", false, "帮助文档, 示例: stress -c 200 -n 10 -u https://www.baidu.com")
	flag.IntVar(&num, "n", 0, "连接数")
	flag.IntVar(&coroutines, "c", 0, "并发数")
	flag.BoolVar(&debug, "d", false, "debug模式")

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	if debug {
		metadata()
	}

	var (
		wg        sync.WaitGroup
		wgReceive sync.WaitGroup
	)

	ch := make(chan *model.ReqResult, 200)
	wgReceive.Add(1)
	go receiveData(ch, &wgReceive)

	for i := 0; i < coroutines; i++ {
		wg.Add(1)
		go stress(url, ch, &wg)
	}

	wg.Wait()
	close(ch)
	time.Sleep(time.Millisecond * 1)
	wgReceive.Wait()

	fmt.Printf("总共花费了%.3f秒", float64(uint64(time.Now().UnixNano()-startTime.UnixNano()))/1e9)
}

// metadata 输出metadata
func metadata() {
	fmt.Printf(`
debug模式: %v
链接:      %s
请求头:    %s
并发数:    %d
并发数量:  %d

`, debug, url, "", coroutines, num)
}

func tableHeader() {
	fmt.Printf("\n")
	fmt.Println(" 耗时│ 并发数│ 成功数│ 失败数│   QPS  │最长耗时│最短耗时│平均耗时│ 错误码")
	fmt.Println("─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────")
}

// stress 压力测试
func stress(url string, ch chan *model.ReqResult, wg *sync.WaitGroup) {
	defer wg.Done()
	var w sync.WaitGroup
	for i := 0; i < num; i++ {
		w.Add(1)
		requestTime := time.Now().UnixNano()
		go process(url, ch, requestTime, &w)
	}
	w.Wait()
}

// process 处理请求
func process(url string, ch chan *model.ReqResult, requestTime int64, wg *sync.WaitGroup) {
	defer wg.Done()
	err := client.Get(url, http.Header{})

	isSuccess := true
	if err != nil {
		isSuccess = false
		if debug {
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

// receiveData 接收数据并将数据输出在控制台上
func receiveData(ch chan *model.ReqResult, wg *sync.WaitGroup) {
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
