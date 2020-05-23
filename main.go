package main

import (
	"flag"
	"fmt"
	"github.com/BUGLAN/stress/client"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	coroutines int
	num        int
	help       bool
	url        string
)

func init() {
	flag.StringVar(&url, "u", "", "链接")
	flag.BoolVar(&help, "h", false, "帮助")
	flag.IntVar(&num, "n", 0, "连接数")
	flag.IntVar(&coroutines, "c", 0, "并发数")

}

func main() {
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup

	startTime := uint64(time.Now().UnixNano())
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	for i := 0; i < coroutines; i++ {
		wg.Add(1)
		go stress(url, &wg)
	}

	wg.Wait()
	fmt.Printf("%.3f秒", float64(uint64(time.Now().UnixNano())-startTime)/1e9)
}

// stress 压力测试
func stress(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < num; i++ {
		time.Sleep(time.Millisecond * 1)
		err := client.Get(url, http.Header{})
		if err != nil {
			log.Println(err)
		}
	}
}
