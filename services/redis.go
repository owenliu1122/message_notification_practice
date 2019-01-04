package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/fpay/foundation-go/log"
	goredis "github.com/go-redis/redis"
)

// MarshalFunc type is an adapter to marshal data for cache to redis.
type MarshalFunc func(interface{}) ([]byte, error)

// UnmarshalFunc type is an adapter to unmarshal cache data.
type UnmarshalFunc func([]byte, interface{}) error

// CacheRedis is redis cache type.
type CacheRedis struct {
	Client    *goredis.Client
	logger    *log.Logger
	Marshal   MarshalFunc
	Unmarshal UnmarshalFunc
}

// NewRedisCli returns a redis cache type client.
func NewRedisCli(logger *log.Logger, url string, marshalfn MarshalFunc, unmarshalfn UnmarshalFunc) (*CacheRedis, error) {

	opt, err := goredis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	fmt.Println("addr is", opt.Addr)
	fmt.Println("db is", opt.DB)
	fmt.Println("password is", opt.Password)

	// Create client as usually.
	c := CacheRedis{
		Client:    goredis.NewClient(opt),
		logger:    logger,
		Marshal:   marshalfn,
		Unmarshal: unmarshalfn,
	}

	return &c, nil
}

// GetClient returns redis-go client.
func (c *CacheRedis) GetClient() (*goredis.Client, error) {
	_, err := c.Client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return c.Client, nil
}

// Close redis cache client.
func (c *CacheRedis) Close() error {
	return c.Client.Close()
}

// Get returns nil error when key does not exist.
func (c *CacheRedis) Get(key string, outVal interface{}) error {

	bytes, err := c.Client.Get(key).Bytes()
	if err != nil {
		c.logger.Errorf("cache: Get %s failed: %s", key, err)
		return err
	}

	if bytes == nil {
		return nil
	}

	if err := c.Unmarshal(bytes, outVal); err != nil {
		c.logger.Errorf("cache: Unmarshal failed: %s", err)
		return err
	}

	return nil
}

// Set key values such as redis set command.
// Zero expiration means the key has no expiration time.
func (c *CacheRedis) Set(key string, inVal interface{}, expiration time.Duration) error {

	bytes, err := c.Marshal(inVal)
	if err != nil {
		c.logger.Errorf("cache: Marshal failed: %s", err)
		return err
	}

	fmt.Printf("Set, bytes: %#v\n", string(bytes))

	_, err = c.Client.Set(key, bytes, expiration).Result()
	if err != nil {
		c.logger.Errorf("cache: Set %s failed: %s", key, err)
	}

	return err
}

// Delete keys from redis cache.
func (c *CacheRedis) Delete(key ...string) error {

	_, err := c.Client.Del(key...).Result()
	if err != nil {
		c.logger.Errorf("cache: Del %s failed: %s", key, err)
	}

	return err
}

// IsExist return key exists value, true or false.
func (c *CacheRedis) IsExist(key string) bool {

	result, err := c.Client.Exists(key).Result()
	if err != nil {
		c.logger.Errorf("cache: IsExist %s failed: %s", key, err)
		return false
	}
	return result > 0
}

// Expire to set key expiration time.
func (c *CacheRedis) Expire(key string, expiration time.Duration) error {
	_, err := c.Client.Expire(key, expiration).Result()
	if err != nil {
		c.logger.Errorf("cache: Expire %s failed: %s", key, err)
	}
	return err
}

// ExpireAt to set key expiration time.
func (c *CacheRedis) ExpireAt(key string, tm time.Time) error {
	_, err := c.Client.ExpireAt(key, tm).Result()
	if err != nil {
		c.logger.Errorf("cache: ExpireAt %s failed: %s", key, err)
	}
	return err
}

// SAdd add a member to current set.
func (c *CacheRedis) SAdd(key string, inVal ...interface{}) error {

	bytesVals := make([]interface{}, 0, len(inVal))

	for i := range inVal {
		bytes, err := c.Marshal(inVal[i])
		if err != nil {
			c.logger.Errorf("cache: Marshal failed: %s", err)
			return err
		}
		bytesVals = append(bytesVals, bytes)
	}

	_, err := c.Client.SAdd(key, bytesVals...).Result()
	if err != nil {
		c.logger.Errorf("cache: SAdd %s failed: %s", key, err)
	}
	return err
}

// SMembers returns all members in current set.
func (c *CacheRedis) SMembers(key string, outSlice interface{}) error {

	valSlice, err := c.Client.SMembers(key).Result()
	if err != nil {
		c.logger.Errorf("cache: SMembers %s failed: %s", key, err)
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
		if e := c.Unmarshal([]byte(s), elem.Addr().Interface()); e != nil {
			e = fmt.Errorf("cache: SMembers index=%d value=%q failed: %s", i, s, e)
			return e
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

//
//type Test struct {
//	Name string `json:"name"`
//	Age  uint8  `json:"age"`
//}
//
//func init() {
//	var (
//		cache *CacheRedis
//		err   error
//	)
//
//	if cache, err = NewRedisCli("redis://localhost:6379/0", json.Marshal, json.Unmarshal); err != nil {
//		panic(err)
//	}
//
//	info := []Test{{"owenjiaxing", 15}, {"ssss", 30}}
//
//	cache.SAdd("owen", info)
//
//	var result []Test
//	cache.SMembers("owen", &result)
//	fmt.Printf("resuslt: %#v\n", result)

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
//}
