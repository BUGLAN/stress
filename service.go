package stress

import (
	"fmt"
	"sync/atomic"
	"net/http"
	"sync"
	"time"

	"github.com/BUGLAN/stress/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ReqResult 单次请求的结构
type ReqResult struct {
	StatusCode  int     // http status code
	GRPCCode    string  // grpc status code
	RequestTime float64 // 处理时间
	IsSuccess   bool    // 是否成功
	err         error   // 错误
}

// StressServer struct
type StressServer struct {
	StartTime    time.Time
	logger       zerolog.Logger
	RequestNum   int // 请求数
	ConnectNum   int // 连接数
	ch           chan *ReqResult
	Url          string
	CoroutineNum int // 同一时间最高并发数
	Handler      Handler
	limiter      chan struct{} // 并发限制
	num          int64
}

// NewStressServer 构造函数
func NewStressServer() *StressServer {
	s := &StressServer{
		logger:  log.With().Str("-", "stress").Caller().Logger(),
		limiter: make(chan struct{}, 255),
		ch:      make(chan *ReqResult, 1000),
	}
	return s
}

// Handler 单次请求和验证的类型
type Handler func(server *StressServer) error

// HttpHandler for standard http handler
func HttpHandler(s *StressServer) error {
	var err error
	httpClient := client.NewHTTPClient()
	requestTime := time.Now().UnixNano()
	err = httpClient.Get(s.Url, http.Header{})
	atomic.AddInt64(&s.num, 1)
	fmt.Println(s.Url)

	isSuccess := true
	if err != nil {
		isSuccess = false
	}
	s.ch <- &ReqResult{
		IsSuccess:   isSuccess,
		StatusCode:  200,
		RequestTime: float64(time.Now().UnixNano() - requestTime),
	}
	return nil
}

// Worker do worker
func (srv *StressServer) Worker(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("start Worker")
	srv.limiter <- struct{}{}
	// TODO need add RequestNum
	if err := srv.Handler(srv); err != nil {
		srv.logger.Warn().Err(err)
	}
	<-srv.limiter
	fmt.Println("end Worker")
}

// Collect the request data
func (srv *StressServer) Collect(wg *sync.WaitGroup) {
	for data := range srv.ch {
		fmt.Println(data)
	}
	fmt.Println("do collect")
	wg.Done()
}

func (srv *StressServer) Start() {
	srv.StartTime = time.Now()
	wg := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go srv.Collect(&wg2)

	// worker
	for i := 0; i < srv.CoroutineNum; i++ {
		wg.Add(1)
		go srv.Worker(&wg)
	}

	wg.Wait()
	close(srv.ch)
	wg2.Wait()

	// 输出最终结果
	fmt.Printf("done")
	fmt.Println(srv.num)
	fmt.Println(float64(srv.num) / time.Now().Sub(srv.StartTime).Seconds())
}


// qps == end_time - start_time
