package main

import (
	"github.com/ijidan/jgo/router"
)

//入口函数
func main() {
	//HTTP开启
	go router.StartHttpServer()
	//阻塞
	ch:=make(chan struct{})
	<-ch

}