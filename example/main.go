package main

import (
	"github.com/BUGLAN/stress"
)

func main() {
	s := stress.NewStressServer()
	s.Url = "https://www.baidu.com"
	s.RequestNum = 1000
	s.CoroutineNum = 100
	s.Handler = stress.StressHTTPHandler
	s.Start()
}
