package jetcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/ijidan/jgo/jgo/jlogger"
	"sync"
	"time"
)

//服务发现
type ServiceDiscovery struct {
	client     *clientv3.Client  //客户
	serverList map[string]string //服务列表
	lock       sync.Mutex
}

//创建客户端
func (sd *ServiceDiscovery) CreateClient(endpoints []string) error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		jlogger.Error("create client error:%s", err.Error())
		return err
	}
	sd.client = client
	return nil
}

//初始化并监听服务
func (sd *ServiceDiscovery) WatchService(prefix string) error {
	rsp, err := sd.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		jlogger.Error("get service error: %s", err.Error())
		return err
	}
	for _, kv := range rsp.Kvs {
		_ = sd.setService(string(kv.Key), string(kv.Value))
	}
	go sd.doWatchService(prefix)
	return nil
}

//监听服务
func (sd *ServiceDiscovery) doWatchService(prefix string) error {
	watchCh := sd.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for watchRsp := range watchCh {
		for _, event := range watchRsp.Events {
			kv := event.Kv
			key := string(kv.Key)
			value := string(kv.Value)
			switch event.Type {
			case mvccpb.PUT:
				_ = sd.setService(key, value)
			case mvccpb.DELETE:
				_ = sd.setService(key, value)
			}
		}
	}
	return nil
}

//设置服务
func (sd *ServiceDiscovery) setService(key string, value string) error {
	sd.lock.Lock()
	defer sd.lock.Unlock()
	sd.serverList[key] = value
	return nil
}

//删除服务
func (sd *ServiceDiscovery) deleteService(key string) error {
	sd.lock.Lock()
	defer sd.lock.Unlock()
	delete(sd.serverList, key)
	return nil
}

//获取服务列表
func (sd *ServiceDiscovery) GetServiceList() map[string]string {
	sd.lock.Lock()
	defer sd.lock.Unlock()
	return sd.serverList
}

//关闭
func (sd *ServiceDiscovery) Close() error {
	return sd.client.Close()
}

//新实例
func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	sd := &ServiceDiscovery{
		serverList: make(map[string]string),
		lock:       sync.Mutex{},
	}
	_ = sd.CreateClient(endpoints)
	return sd
}
