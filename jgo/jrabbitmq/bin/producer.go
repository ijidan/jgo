package main

import (
	"github.com/ijidan/jgo/jgo/jrabbitmq"
	"strconv"
)

func main()  {
	jr:=jrabbitmq.NewJRabbit()
	//发送消息
	cnt:=100
	for i:=1;i<=cnt;i++{
		_ = jr.PublishMessage("'", "jidan_q_1", "hello jidan "+strconv.Itoa(i), true, false)
	}
}