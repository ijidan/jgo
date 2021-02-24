package main

import (
	"jgo/protected/jgo/jetcd"
	"log"
	"time"
)

func main()  {
	//etcd:=jetcd.NewJEtcd("http://www.jgo.com")
	//rsp:=etcd.Put("name","jidan")
	//log.Println(rsp)
	//jlogger.Info("done...")

	var endpoints = []string{jetcd.EtcdUrl}
	sd:=jetcd.NewServiceDiscovery(endpoints)
	defer func() {
		_=sd.Close()
	}()
	_ = sd.WatchService("/service/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(sd.GetServiceList())
		}
	}
}
