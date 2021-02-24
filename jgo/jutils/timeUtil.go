package jutils

import (
	"strconv"
	"time"
)

//时间工具
type timeUtil struct {

}

//获取当前时间戳
func (t *timeUtil) GetCurrentTimeStamp()string {
	timeStamp := time.Now().Unix()
	timeStr := strconv.FormatInt(timeStamp, 10)
	return timeStr
}

//获取当前时间
func (t *timeUtil) GetCurrentTime()string  {
	timeUnix := time.Now().Unix()
	time_ := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	return time_
}

//获取实例
func NewTimeUtil() *timeUtil  {
	t:=&timeUtil{}
	return t
}