package stress

import (
	"fmt"
	"github.com/BUGLAN/stress/client"
	"github.com/BUGLAN/stress/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
	"time"
)

type StressServer struct {
	logger       zerolog.Logger
	RequestNum   int // 请求数
	ch           chan *model.ReqResult
	Url          string
	CoroutineNum int // 同一时间最高并发数
	Handler      Handler
	limiter      chan struct{} // 并发限制
}

func NewStressServer() *StressServer {
	s := &StressServer{
		logger:  log.With().Str("-", "stress").Caller().Logger(),
		limiter: make(chan struct{}, 255),
	}
	return s
}

type Handler func(server *StressServer) error

// DoStress do stress request
func (srv *StressServer) DoHandlers(handlers ...Handler) {
	var err error
	for _, f := range handlers {
		if err = f(srv); err != nil {
			srv.logger.Warn().Err(err)
		}
	}
}

func StressHTTPHandler(s *StressServer) error {
	var err error
	httpClient := client.NewHTTPClient()
	requestTime := time.Now().UnixNano()
	err = httpClient.Get(s.Url, http.Header{})

	isSuccess := true
	if err != nil {
		isSuccess = false
		if model.Debug {
			fmt.Printf("err: %s\n", err.Error())
		}
	}
	// 每个请求都开协程的话, 消耗过大, 反而不利于并发请求
	s.ch <- &model.ReqResult{
		IsSuccess:   isSuccess,
		StatusCode:  200,
		RequestTime: float64(time.Now().UnixNano() - requestTime),
	}
	return nil
}

// Worker do worker
func (srv *StressServer) Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	srv.limiter <- struct{}{}
	if err := srv.Handler(srv); err != nil {
		srv.logger.Warn().Err(err)
	}
	<-srv.limiter
}

// Collect the request data
func (srv *StressServer) Collect(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("do collect")
}

func (srv *StressServer) Start() {
	wg := sync.WaitGroup{}
	ch := make(chan *model.ReqResult, 1000)
	wg.Add(1)
	go srv.Collect(&wg)

	// worker
	for i := 0; i < srv.RequestNum; i++ {
		wg.Add(1)
		go srv.Worker(&wg)
	}

	wg.Wait()
	close(ch)
	fmt.Printf("done")
}
