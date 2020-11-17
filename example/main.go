package main

import (
	"github.com/BUGLAN/stress"
)

func main() {
	s := stress.NewStressServer()
	s.Url = "https://www.baidu.com"
	s.RequestNum = 100
	s.CoroutineNum = 100
	s.Handler = stress.HttpHandler
	s.Start()
}
