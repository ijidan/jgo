package router

import (
	"fmt"
	"github.com/ijidan/jgo/controller"
	"github.com/ijidan/jgo/jgo/jlogger"
	"github.com/ijidan/jgo/jgo/jrouter"
	"net/http"
)

const HttpHost = "127.0.0.1"
const HttpPort = int64(8080)

//注册n
func Registry() {
	//控制器
	index := controller.IndexController{}
	user := controller.UserController{}
	chat := controller.ChatController{}
	admin := controller.AdminIndexController{}
	client := controller.ClientController{}
	server := controller.ServerController{}
	service:=controller.ServiceController{}

	//路由实例
	ins := jrouter.NewJRouter()

	//用户前台
	ins.Group("front", func() {
		//用户前台
		ins.Any("/", index.Index)
		ins.Any("/call", index.Call)
		ins.Any("/user/reg", user.Reg)
		ins.Any("/user/login", user.Login)
		ins.Any("/user/info", user.Info)
		ins.Any("/user/logout", user.Logout)
		ins.Any("/chat/index", chat.Index)
	}).Add(LogMiddleware, BridgeAccountCookieSyncMiddleware)
	
	//管理后台
	ins.Group("backend", func() {
		ins.Any("/admin", admin.Index)
		ins.Any("/admin/defaultIdx", admin.DefaultIdx)
		ins.Any("/admin/serverList", admin.ServerList)
		ins.Any("/admin/clientList", admin.ClientList)
	}).Add(LogMiddleware)

	//API接口
	ins.Group("api", func() {
		ins.Any("/api/client/getAll", client.GetAll)
		ins.Any("/api/client/sendMessage", client.SendMessage)
		ins.Any("/api/client/kickOff", client.KickOff)

		ins.Any("/api/server/getAll", server.GetAll)
		ins.Any("/api/server/batchSendMessage", server.BatchSendMessage)
		ins.Any("/api/server/closeServer", server.CloseServer)
	}).Add(verifyTokenMiddleware)
	
	//Service
	ins.Group("service", func() {
		ins.Any("/service/consulCheck", service.ConsulCheck)
	})
}

//HTTP开启
const staticPath = "protected/static/"
const staticUrlPrefix = "/static/"


//new mux
func newMux() http.Handler {
	//处理
	mux := http.NewServeMux()
	ins := jrouter.NewJRouter()
	allRoute := ins.GetAll() //获取所有路由
	for k, v := range allRoute {
		mux.Handle(k, v)
	}
	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle(staticUrlPrefix, http.StripPrefix(staticUrlPrefix, fs))
	return mux
}

//开启HTTP服务
func StartHttpServer() {
	//HTTP服务
	Registry()
	mux:=newMux()
	address := fmt.Sprintf("%s:%d", HttpHost, HttpPort)
	err := http.ListenAndServe(address, mux)
	if err != nil {
		jlogger.Error(err.Error())
	}
}
