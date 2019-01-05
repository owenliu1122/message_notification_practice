package services

import (
	"encoding/json"

	"github.com/owenliu1122/notice"
)

type Job notice.Job

// Queue 消息队列名
func (j *Job) Queue() string {
	return j.Q
}

// Delay 消息队列延时时间
func (j *Job) Delay() int {
	return j.D
}

// Marshal 序列化将要存入队列的消息
func (j *Job) Marshal() ([]byte, error) {
	return json.Marshal(j.Message)
}
