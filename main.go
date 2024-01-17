package main

import (
	"fmt"
	"logv1/conf"
	"logv1/etcd"
	"logv1/kafka"
	"logv1/taillog"
	"time"
)

var (
	cfg conf.AppConf
)

func main() {
	cfg = conf.Appconf
	// 1.初始化kafka
	err := kafka.Init([]string{cfg.Kafka.Addr}, cfg.Kafka.ChanMaxSize)
	if err != nil {
		fmt.Printf("init kafka failed ,err:%v\n", err)
		return
	}
	fmt.Println("init kafka success")

	// 2. 初始化etcd
	err = etcd.Init(cfg.Etcd.Address, time.Duration(cfg.Etcd.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("init etcd success.")

	// 拉取独自的配置
	// ipStr, err := utils.GetOutboundIP(cfg.Center.Addr)
	// if err != nil {
	// 	panic(err) // 获取不到，直接panic
	// }
	// etcdConfKey := fmt.Sprintf(cfg.Etcd.Key, ipStr)
	etcdConfKey := cfg.Etcd.Key
	// 2.1 从etcd中获取日志收集项的配置信息
	logEntryConf, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Printf("get conf from etcd failed,err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntryConf)
	// 显示etcd有哪些配置
	for index, value := range logEntryConf {
		fmt.Printf("index:%v value:%v\n", index, value)
	}

	// 3. 初始化tail任务，监听路径下对应的log文件信息并将新信息发送到kafka
	taillog.Init(logEntryConf)

	// 4. 监听新配置
	newConfChan := taillog.NewConfChan()         // 对外暴露的更新conf通道
	go etcd.WatchConf(cfg.Etcd.Key, newConfChan) // 哨兵发现最新的配置信息会通知上面的那个通道

	// hang on
	run := make(chan struct{})
	<-run

}
