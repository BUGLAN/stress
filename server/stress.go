package server

import (
	"sync"

	"github.com/BUGLAN/stress/model"
)

type Stress interface {
	Worker(url string, ch chan *model.ReqResult, wg *sync.WaitGroup)
	Receiver(coroutineNum uint64, ch chan *model.ReqResult, wg *sync.WaitGroup)
}

type Server struct{}

func NewServer() *Server {
	return &Server{}
}
