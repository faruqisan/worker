package cache

import (
	"github.com/go-redis/redis"
)

type (
	// Engine struct ..
	Engine struct {
		*redis.Client
	}
)

// New function return engine struct with setuped redis client
func New() Engine {
	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return Engine{c}
}
