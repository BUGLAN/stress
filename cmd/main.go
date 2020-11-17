package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/BUGLAN/stress/model"
	"github.com/BUGLAN/stress/server"
)

// application entrance
func main() {
	// set runtime process
	runtime.GOMAXPROCS(1)

	// set start time
	startTime := time.Now()

	// set flag vars
	flag.StringVar(&model.Url, "u", "", "url 链接 https://www.baidu.com")
	flag.BoolVar(&model.Help, "h", false, "帮助文档, 示例: stress -c 200 -n 10 -u https://www.baidu.com")
	flag.IntVar(&model.Num, "n", 0, "连接数")
	flag.Uint64Var(&model.CoroutineNum, "c", 0, "并发数")
	flag.BoolVar(&model.Debug, "d", false, "debug模式")
	flag.Parse()

	// hint
	if model.Help {
		flag.Usage()
		return
	}

	if model.Debug {
		metadata()
	}

	var (
		wg        sync.WaitGroup
		wgReceive sync.WaitGroup
	)

	// create server
	srv := server.NewServer()
	ch := make(chan *model.ReqResult, 1000)
	wgReceive.Add(1)
	go srv.Receiver(model.CoroutineNum, ch, &wgReceive)

	// do worker
	for i := 0; i < int(model.CoroutineNum); i++ {
		wg.Add(1)
		go srv.Worker(model.Url, ch, &wg)
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

`, model.Debug, model.Url, "", model.CoroutineNum, model.Num)
}
