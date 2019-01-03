package notice

import "time"

// Cache is a redis cache interface
type Cache interface {
	Get(key string, outVal interface{}) error
	Set(key string, inVal interface{}, timeout time.Duration) error
	Delete(key ...string) error
	IsExist(key string) bool
	Expire(key string, expiration time.Duration) error
	ExpireAt(key string, tm time.Time) error
	Close() error
	SAdd(key string, inVal ...interface{}) error
	SMembers(key string, outSlice interface{}) error
}
