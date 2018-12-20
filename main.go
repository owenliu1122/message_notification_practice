//package main
//
//import (
//	"fmt"
//	"message_notification_practice/cmd"
//	"message_notification_practice/redis"
//	_ "message_notification_practice/mysql"
//)
//
//func main(){
//	cmd.Execute()
//	redisCli := redis.GetRedisCli()
//	if err := redisCli.Set("aaa", "999",0).Err(); err == nil {
//		val, _ := redisCli.Get("aaa").Result()
//		fmt.Printf("aaa= %s\n", val)
//	}
//
//
//}

package main

import (
	"bytes"
	"message_notification_practice/cmd"
	_ "message_notification_practice/mysql"
)

func main() {
	cmd.Execute()
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}
