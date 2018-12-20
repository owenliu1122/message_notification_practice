package redis

import (
	"fmt"
	goredis "github.com/go-redis/redis"
	"log"
)

var redisdb *goredis.Client

func init() {

	redisdb = goredis.NewClient(&goredis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisdb.Ping().Result()
	log.Debug(pong, err)
}

func GetRedisCli() *goredis.Client {

	pong, err := redisdb.Ping().Result()
	if err != nil {
		log.Println(pong, err)

		redisdb.Close()

		redisdb = goredis.NewClient(&goredis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}

	return redisdb
}
