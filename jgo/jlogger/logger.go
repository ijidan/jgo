package jlogger

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
	"time"
)

//获取code func
func getData() (string, string, string) {
	pc, codePath, codeLine, ok := runtime.Caller(2)
	var code_ string
	var func_ string
	if !ok {
		code_ = "-"
		func_ = "-"
	} else {
		code_ = fmt.Sprintf("%s:%d", codePath, codeLine)
		func_ = runtime.FuncForPC(pc).Name()
	}
	timeUnix := time.Now().Unix()
	time_ := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	return time_, code_, func_
}

//信息
func Info(format string, a ...interface{}) {
	var message string
	if a == nil {
		message = format
	} else {
		message = fmt.Sprintf(format, a)
	}
	level_ := "Info"
	time_, code_, func_ := getData()
	msg := fmt.Sprintf("%s [%s] %s : %s %s ", time_, level_, message, code_, func_)
	color.White(msg)
}

//注意
func Notice(format string, a ...interface{}) {
	var message string
	if a == nil {
		message = format
	} else {
		message = fmt.Sprintf(format, a)
	}
	level_ := "Notice"
	time_, code_, func_ := getData()
	msg := fmt.Sprintf("%s [%s] %s : %s %s ", time_, level_, message, code_, func_)
	color.Green(msg)
}

//警告
func Warning(format string, a ...interface{}) {
	var message string
	if a == nil {
		message = format
	} else {
		message = fmt.Sprintf(format, a)
	}
	level_ := "Warning"
	time_, code_, func_ := getData()
	msg := fmt.Sprintf("%s [%s] %s : %s %s ", time_, level_, message, code_, func_)
	color.Yellow(msg)
}

//错误
func Error(format string, a ...interface{}) {
	var message string
	if a == nil {
		message = format
	} else {
		message = fmt.Sprintf(format, a)
	}
	level_ := "Error"
	time_, code_, func_ := getData()
	msg := fmt.Sprintf("%s [%s] %s : %s %s ", time_, level_, message, code_, func_)
	color.Red(msg)
}
