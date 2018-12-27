package redis

import (
	"fmt"
	goredis "github.com/go-redis/redis"
	log "gopkg.in/cihub/seelog.v2"
	"time"
)

type Cache interface {
	Get(key string, outVal interface{}) error
	Set(key string, inVal interface{}, timeout time.Duration) error
	Delete(key string) error
	IsExist(key string) bool
	Close() error
}

type MarshalFunc func(interface{}) ([]byte, error)
type UnmarshalFunc func([]byte, interface{}) error

type CacheRedis struct {
	Client    *goredis.Client
	Marshal   MarshalFunc
	Unmarshal UnmarshalFunc
}

func NewRedisCli(url string, marshalfn MarshalFunc, unmarshalfn UnmarshalFunc) (*CacheRedis, error) {

	opt, err := goredis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	fmt.Println("addr is", opt.Addr)
	fmt.Println("db is", opt.DB)
	fmt.Println("password is", opt.Password)

	// Create client as usually.
	cache := CacheRedis{
		Client:    goredis.NewClient(opt),
		Marshal:   marshalfn,
		Unmarshal: unmarshalfn,
	}

	return &cache, nil
}

func (cache *CacheRedis) GetClient() (*goredis.Client, error) {
	_, err := cache.Client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return cache.Client, nil
}

func (cache *CacheRedis) Close() error {
	return cache.Client.Close()
}

func (cache *CacheRedis) Get(key string, outVal interface{}) error {

	bytes, err := cache.Client.Get(key).Bytes()
	if err != nil {
		log.Errorf("cache: Get %s failed: %s", key, err)
		return err
	}

	if bytes == nil {
		return nil
	}

	if err := cache.Unmarshal(bytes, outVal); err != nil {
		log.Errorf("cache: Unmarshal failed: %s", err)
		return err
	}

	return nil
}

func (cache *CacheRedis) Set(key string, inVal interface{}, expiration time.Duration) error {

	bytes, err := cache.Marshal(inVal)
	if err != nil {
		log.Errorf("cache: Marshal failed: %s", err)
		return err
	}

	fmt.Printf("Set, bytes: %#v\n", string(bytes))

	_, err = cache.Client.Set(key, bytes, expiration).Result()
	if err != nil {
		log.Errorf("cache: Set %s failed: %s", key, err)
	}

	return err
}

func (cache *CacheRedis) Delete(key string) error {

	_, err := cache.Client.Del(key).Result()
	if err != nil {
		log.Errorf("cache: Del %s failed: %s", key, err)
	}

	return err
}

func (cache *CacheRedis) IsExist(key string) bool {

	result, err := cache.Client.Exists(key).Result()
	if err != nil {
		log.Errorf("cache: IsExist %s failed: %s", key, err)
		return false
	}
	return result > 0
}

//type Test struct {
//	Name string		`json:"name"`
//	Age  uint8		`json:"age"`
//}
//
//func init() {
//	var (
//		cache *CacheRedis
//		err error
//	)
//
//	if cache, err = NewRedisCli("redis://localhost:6379/0", json.Marshal, json.Unmarshal); err!=nil {
//		panic(err)
//	}
//
//	info:=&Test{"owenjiaxing",15}
//	cache.Set("myinfo", info, 10*time.Second)
//
//	var infoRet Test
//
//	cache.Get("myinfo", &infoRet)
//	fmt.Printf("infoRet: %#v\n", infoRet)
//	time.Sleep(5*time.Second)
//	fmt.Printf("myinfo exists: %v\n", cache.Delete("myinfo"))
//}
