package mq

import (
	"fmt"
	"log"
)

type MQCfg struct {
	URL      string `json:"hosts"`
	Queue    string `json:"queue"`
	Exchange string `json:"exchange"`
}

type MQInfo struct {
	Cfg     MQCfg
	MsgChan chan interface{} // 消息队列
	//ExtChan   chan bool         // 协程退出标记信号
	//WaitGroup *sync.WaitGroup   // 同步组
}

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}
