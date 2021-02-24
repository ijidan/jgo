package jetcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/ijidan/jgo/jgo/jlogger"
	"google.golang.org/grpc/grpclog"
	"os"
	"time"
)

//服务注册
type ServiceRegister struct {
	client             *clientv3.Client //客户端
	leaseId            clientv3.LeaseID //租约ID
	leaseKeepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key                string
	value              string
}

//创建客户端
func (sr *ServiceRegister) CreateClient(endpoints []string) error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		jlogger.Error("create client error:%s", err.Error())
		return err
	}
	sr.client = client
	return nil
}

//生成租约
func (sr *ServiceRegister) GrantLease(ttl int64) (clientv3.LeaseID, error) {
	rsp, err := sr.client.Grant(context.Background(), ttl)
	if err != nil {
		jlogger.Error("grant lease error: %s", err.Error())
		return 0, err
	}
	return rsp.ID, nil
}

//撤销租约
func (sr *ServiceRegister) RevokeLease(leaseId clientv3.LeaseID) error {
	_, err := sr.client.Revoke(context.Background(), leaseId)
	if err != nil {
		jlogger.Error("revoke lease error:%s", err.Error())
		return err
	}
	_ = sr.Close()
	return nil
}

//续约
func (sr *ServiceRegister) KeepAliveLease(leaseId clientv3.LeaseID) error {
	kaChan, err := sr.client.KeepAlive(context.Background(), leaseId)
	if err != nil {
		jlogger.Error("keep alive lease error:%s", err.Error())
		return err
	}
	sr.leaseKeepAliveChan = kaChan
	return nil
}

//设置KV
func (sr *ServiceRegister) PutKeyWithLease(leaseId clientv3.LeaseID, key string, value string) error {
	opts := clientv3.WithLease(leaseId)
	_, err := sr.client.Put(context.Background(), key, value, opts)
	if err != nil {
		jlogger.Error("put key error: %s", err.Error())
		return err
	}
	return nil
}

//监听
func (sr *ServiceRegister) ListLeaseKeepAliveChan() error{
	for rsp := range sr.leaseKeepAliveChan {
		jlogger.Info("keep lease alive success:%d", rsp.ID)
	}
	return nil
}

//关闭
func (sr *ServiceRegister) Close() error {
	return sr.client.Close()
}

//新实例
func NewServiceRegister(endpoints []string, ttl int64, key string, value string) *ServiceRegister {
	clientv3.SetLogger(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))
	sr := &ServiceRegister{}
	_= sr.CreateClient(endpoints)
	leaseId, _ := sr.GrantLease(ttl)
	_ = sr.PutKeyWithLease(leaseId, key, value)
	return sr
}
