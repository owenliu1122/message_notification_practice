package redis

import (
	"encoding/json"
	"fmt"
	goredis "github.com/go-redis/redis"
	log "gopkg.in/cihub/seelog.v2"
	"reflect"
	"time"
)

type Cache interface {
	Get(key string, outVal interface{}) error
	Set(key string, inVal interface{}, timeout time.Duration) error
	Delete(key string) error
	IsExist(key string) bool
	Expire(key string, expiration time.Duration) error
	ExpireAt(key string, tm time.Time) error
	Close() error
	// Set
	SAdd(key string, inVal ...interface{}) error
	SMembers(key string, outSlice interface{}) error
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

func (cache *CacheRedis) Expire(key string, expiration time.Duration) error {
	_, err := cache.Client.Expire(key, expiration).Result()
	if err != nil {
		log.Errorf("cache: Expire %s failed: %s", key, err)
	}
	return err
}

func (cache *CacheRedis) ExpireAt(key string, tm time.Time) error {
	_, err := cache.Client.ExpireAt(key, tm).Result()
	if err != nil {
		log.Errorf("cache: ExpireAt %s failed: %s", key, err)
	}
	return err
}

/* Set */

func (cache *CacheRedis) SAdd(key string, inVal ...interface{}) error {

	bytesVals := make([]interface{}, 0, len(inVal))

	for i := range inVal {
		bytes, err := cache.Marshal(inVal[i])
		if err != nil {
			log.Errorf("cache: Marshal failed: %s", err)
			return err
		}
		bytesVals = append(bytesVals, bytes)
		fmt.Printf("[%d], bytes: %s\n", string(bytes))
	}

	_, err := cache.Client.SAdd(key, bytesVals...).Result()
	if err != nil {
		log.Errorf("cache: SAdd %s failed: %s", key, err)
	}
	return err
}

func (cache *CacheRedis) SMembers(key string, outSlice interface{}) error {

	valSlice, err := cache.Client.SMembers(key).Result()
	if err != nil {
		log.Errorf("cache: SMembers %s failed: %s", key, err)
		return err
	}

	if len(valSlice) == 0 {
		return nil
	}

	v := reflect.ValueOf(outSlice)
	if !v.IsValid() {
		return fmt.Errorf("cache: SMembers(nil)")
	}
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("cache: SMembers(non-pointer %T)", outSlice)
	}
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("cache: SMembers(non-slice %T)", outSlice)
	}

	next := makeSliceNextElemFunc(v)
	for i, s := range valSlice {
		elem := next()
		if err := cache.Unmarshal([]byte(s), elem.Addr().Interface()); err != nil {
			err = fmt.Errorf("cache: SMembers index=%d value=%q failed: %s", i, s, err)
			return err
		}
	}

	return err
}

func makeSliceNextElemFunc(v reflect.Value) func() reflect.Value {
	elemType := v.Type().Elem()

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
		return func() reflect.Value {
			if v.Len() < v.Cap() {
				v.Set(v.Slice(0, v.Len()+1))
				elem := v.Index(v.Len() - 1)
				if elem.IsNil() {
					elem.Set(reflect.New(elemType))
				}
				return elem.Elem()
			}

			elem := reflect.New(elemType)
			v.Set(reflect.Append(v, elem))
			return elem.Elem()
		}
	}

	zero := reflect.Zero(elemType)
	return func() reflect.Value {
		if v.Len() < v.Cap() {
			v.Set(v.Slice(0, v.Len()+1))
			return v.Index(v.Len() - 1)
		}

		v.Set(reflect.Append(v, zero))
		return v.Index(v.Len() - 1)
	}
}

type Test struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

func init() {
	var (
		cache *CacheRedis
		err   error
	)

	if cache, err = NewRedisCli("redis://localhost:6379/0", json.Marshal, json.Unmarshal); err != nil {
		panic(err)
	}

	info := []Test{{"owenjiaxing", 15}, {"ssss", 30}}

	cache.SAdd("owen", info[0], info[1])

	var result []Test
	cache.SMembers("owen", &result)
	fmt.Printf("resuslt: %#v\n", result)

	//cache.ExpireAt("owen", time.Now())
	//cache.ExpireAt("owen", time.Now())

	//cache.Set("myinfo_1", info, 60*time.Second)
	//info.Age = 100
	//cache.Set("myinfo_2", info, 60*time.Second)
	//
	//var infoRet []Test
	//
	//cache.Get("myinfo*", &infoRet)
	//fmt.Printf("infoRet: %#v\n", infoRet)
	//time.Sleep(5*time.Second)
	//fmt.Printf("myinfo exists: %v\n", cache.Delete("myinfo"))
}
