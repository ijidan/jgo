package jrouter

import (
	"github.com/ijidan/jgo/jgo/jcontext"
	"log"
)

//路由检测中间件
func DefaultCheckRouterMiddleware(c *jcontext.Context) (next bool, err error) {
	log.Println("默认路由检测中间件")
	return true, nil
}
