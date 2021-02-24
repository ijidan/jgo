package jrouter

import (
	"context"
	"github.com/ijidan/jgo/controller"
	"github.com/ijidan/jgo/jgo/jcontext"
	"github.com/ijidan/jgo/jgo/jjaeger"
	"github.com/ijidan/jgo/jgo/jlogger"
	"github.com/ijidan/jgo/jgo/jutils"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"sync"
)

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodPut    = "PUT"
	methodDelete = "DELETE"
	methodAny    = "ANY"
)

//中间件
type JMiddlewareHandler func(c *jcontext.Context) (next bool, err error)

//处理函数
type JHandlerFunc func(c *jcontext.Context)

//分组函数
type JGroup func()

//路由
type JRouter struct {
	JRouterHandlerMap    map[string]JHandlerFunc
	JRouterMethodMap     map[string]string
	JRouterGroupMap      map[string]string
	JRouterMiddlewareMap map[string]jutils.ArrayStack

	//分组相关
	JGroupName      string
	JGroupRouterMap map[string][]string
}

//GET请求
func (jr *JRouter) Get(match string, f JHandlerFunc, handlerList ...JMiddlewareHandler) {
	if f != nil {
		jr.JRouterHandlerMap[match] = f
		jr.JRouterMethodMap[match] = methodGet
		if len(jr.JGroupName) > 0 {
			jr.JRouterGroupMap[match] = jr.JGroupName
			jr.JGroupRouterMap[jr.JGroupName] = append(jr.JGroupRouterMap[jr.JGroupName], match)
		}
	}
	if len(handlerList) > 0 {
		jr.addMiddleware(match, handlerList)
	}
}

//POST请求
func (jr *JRouter) Post(match string, f JHandlerFunc, handlerList ...JMiddlewareHandler) {
	if f != nil {
		jr.JRouterHandlerMap[match] = f
		jr.JRouterMethodMap[match] = methodPost
		if len(jr.JGroupName) > 0 {
			jr.JRouterGroupMap[match] = jr.JGroupName
			jr.JGroupRouterMap[jr.JGroupName] = append(jr.JGroupRouterMap[jr.JGroupName], match)
		}
	}
	if len(handlerList) > 0 {
		jr.addMiddleware(match, handlerList)
	}
}

//PUT请求
func (jr *JRouter) Put(match string, f JHandlerFunc, handlerList ...JMiddlewareHandler) {
	if f != nil {
		jr.JRouterHandlerMap[match] = f
		jr.JRouterMethodMap[match] = methodPut
		if len(jr.JGroupName) > 0 {
			jr.JRouterGroupMap[match] = jr.JGroupName
			jr.JGroupRouterMap[jr.JGroupName] = append(jr.JGroupRouterMap[jr.JGroupName], match)
		}
	}
	if len(handlerList) > 0 {
		jr.addMiddleware(match, handlerList)
	}
}

//delete请求
func (jr *JRouter) Delete(match string, f JHandlerFunc, handlerList ...JMiddlewareHandler) {
	if f != nil {
		jr.JRouterHandlerMap[match] = f
		jr.JRouterMethodMap[match] = methodDelete
		if len(jr.JGroupName) > 0 {
			jr.JRouterGroupMap[match] = jr.JGroupName
			jr.JGroupRouterMap[jr.JGroupName] = append(jr.JGroupRouterMap[jr.JGroupName], match)
		}
	}
	if len(handlerList) > 0 {
		jr.addMiddleware(match, handlerList)
	}
}

//添加路由
func (jr *JRouter) Any(match string, f JHandlerFunc, handlerList ...JMiddlewareHandler) {
	if f != nil {
		jr.JRouterHandlerMap[match] = f
		jr.JRouterMethodMap[match] = methodAny
		if len(jr.JGroupName) > 0 {
			jr.JRouterGroupMap[match] = jr.JGroupName
			jr.JGroupRouterMap[jr.JGroupName] = append(jr.JGroupRouterMap[jr.JGroupName], match)
		}
	}
	if len(handlerList) > 0 {
		jr.addMiddleware(match, handlerList)
	}
}

//分组
func (jr *JRouter) Group(groupName string, f JGroup) *JRouter {
	jr.JGroupName = groupName
	f()
	return jr
}

//添加中间件
func (jr *JRouter) Add(handlerList ...JMiddlewareHandler) {
	if len(jr.JGroupName) > 0 && len(handlerList) > 0 {
		//添加中间件
		if matches, ok := jr.JGroupRouterMap[jr.JGroupName]; ok {
			for _, match := range matches {
				jr.addMiddleware(match, handlerList)
			}
		}
	}
	jr.JGroupName = ""
}

//添加中间件
func (jr *JRouter) addMiddleware(match string, handlerList []JMiddlewareHandler) {
	stack := jr.computeStack(match, handlerList)
	jr.JRouterMiddlewareMap[match] = *stack
}

//计算堆栈
func (jr *JRouter) computeStack(name string, handlerList []JMiddlewareHandler) *jutils.ArrayStack {
	stack := jutils.NewArrayStack(name)
	for _, handler := range handlerList {
		stack.Push(handler)
	}
	return stack
}

//执行中间件
func (jr *JRouter) ExecMiddleware(match string, c *jcontext.Context) (next bool, err error) {
	//添加默认中间件
	jr.Any(match, nil, DefaultCheckRouterMiddleware)
	stack := jr.JRouterMiddlewareMap[match]
	//执行自定义中间件
	for {
		topEle := stack.Pop()
		if topEle == nil {
			break
		}
		handler := topEle.(JMiddlewareHandler)
		next, err = handler(c)
		if !next || err != nil {
			return next, err
		}
	}
	return true, nil
}

//获取所有路由
func (jr *JRouter) GetAll() map[string]http.HandlerFunc {
	routerMap := make(map[string]http.HandlerFunc)
	for url, _ := range jr.JRouterHandlerMap {
		m := func(routerHandlerMap map[string]JHandlerFunc, routerMethodMap map[string]string) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				//recover handler
				defer func() {
					ctx := jcontext.NewContext(w, r, nil)
					if r := recover(); r != nil {
						jlogger.Error("recover received：", r)
						sys := controller.SystemController{}
						sys.ErrorInner(ctx)
					}
				}()
				path := r.URL.Path
				method := r.Method
				//遍历
				if _, ok := routerMethodMap[path]; ok {
					assertMethod := routerMethodMap[path]
					if assertMethod == methodAny || assertMethod == method {
						//执行默认中间件
						if _, ok := routerHandlerMap[path]; ok {
							//jaeger
							tracer, closer,span := jjaeger.NewJJaeger("jaeger-service-call","jaeger-span-call-root")
							opentracing.SetGlobalTracer(tracer)
							bgCtx:=context.Background()
							tracingCtx := opentracing.ContextWithSpan(bgCtx, span)
							ctx := jcontext.NewContext(w, r, tracingCtx)
							//执行中间件
							_, err := jr.ExecMiddleware(path, ctx)
							if err != nil {
								_, _ = w.Write([]byte("中间件错误：" + err.Error()))
							} else {
								//控制器处理
								handlerFunc := routerHandlerMap[path]
								handlerFunc(ctx)
							}
							span.Finish()
							err=closer.Close()
							if err != nil {
								_, _ = w.Write([]byte("closer错误：" + err.Error()))
							}
						}
					} else {
						_, _ = w.Write([]byte("请求方式错误，支持：" + assertMethod))
					}
				}

			}
		}(jr.JRouterHandlerMap, jr.JRouterMethodMap)
		routerMap[url] = m
	}
	return routerMap
}

//实例
var instance *JRouter
var once sync.Once

//获取单例
func NewJRouter() *JRouter {
	once.Do(func() {
		instance = &JRouter{
			JRouterHandlerMap:    make(map[string]JHandlerFunc),
			JRouterMethodMap:     map[string]string{},
			JRouterGroupMap:      make(map[string]string),
			JGroupName:           "",
			JGroupRouterMap:      make(map[string][]string),
			JRouterMiddlewareMap: make(map[string]jutils.ArrayStack),
		}
	})
	return instance
}
