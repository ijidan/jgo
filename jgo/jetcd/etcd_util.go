package jetcd

import (
	"encoding/json"
	"github.com/ijidan/jnet/jnet"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

//etcd 工具
type JEtcd struct {
	jnet.BaseService
	host  string
	param map[string]interface{}
}

//设置数据
func (et *JEtcd) Put(key string, value string) jnet.Response {
	url := et.host + "/v3/kv/put"
	reqParam := map[string]interface{}{"key": key, "value": value}
	et.SetUseProxy(true)
	response := et.SendPostRequest(url, reqParam)
	rsp := et.computeResponse(response)
	return rsp
}

//获取数据
func (et *JEtcd) Get(key string) jnet.Response {
	url := et.host + "/v3/kv/range"
	reqParam := map[string]interface{}{"key": key}
	et.SetUseProxy(true)
	response := et.SendPostRequest(url, reqParam)
	rsp := et.computeResponse(response)
	return rsp
}

//根据前缀获取数据
func (et *JEtcd) GetWithPrefix(prefix string, rangeEnd string) jnet.Response {
	url := et.host + "/v3/kv/range"
	reqParam := map[string]interface{}{"key": prefix, "range_end": rangeEnd}
	et.SetUseProxy(true)
	response := et.SendPostRequest(url, reqParam)
	rsp := et.computeResponse(response)

	return rsp
}

//响应
func (et *JEtcd) computeResponse(rsp jnet.Response) jnet.Response {
	isFail := rsp.Fail()
	if isFail {
		return jnet.NewResponse(rsp.GetCode(), rsp.GetMessage(), rsp.GetData(), rsp.GetPrompt())
	}
	body := rsp.GetData()
	bodyString := body.(string)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(bodyString), &m); err != nil {
		return jnet.NewResponse(jnet.JsonParseFail, err.Error(), nil, "")
	}
	return jnet.NewResponse(jnet.Success, "", m["data"], "")
}

var instanceMap map[string]*JEtcd
var lock sync.Mutex

//获取实例
func NewJEtcd(host string) *JEtcd {
	lock.Lock()
	defer lock.Unlock()
	if instanceMap == nil {
		instanceMap = make(map[string]*JEtcd)
	}
	if _, ok := instanceMap[host]; !ok {
		instance := &JEtcd{
			host:  host,
			param: make(map[string]interface{}),
		}
		instanceMap[host] = instance
	}
	instance := instanceMap[host]
	instance.param = make(map[string]interface{})
	return instanceMap[host]

}

//获取服务列表
func GetServiceList() {
	var endpoints = []string{EtcdUrl}
	sd := NewServiceDiscovery(endpoints)
	defer func() {
		_ = sd.Close()
	}()
	_ = sd.WatchService("/service/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(sd.GetServiceList())
		}
	}
}
