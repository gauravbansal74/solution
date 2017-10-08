package queue

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	Host   string
	Port   int
	DB     int
	Client *redis.Client
}

var (
	redisClient RedisClient
)

func (r *RedisClient) urlFor() string {
	return fmt.Sprintf(
		"%s:%d",
		r.Host,
		r.Port,
	)
}

func (r *RedisClient) Init() {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.urlFor(),
		Password: "",
		DB:       r.DB,
	})
	redisClient = *r
}

func ReadRedisClient() RedisClient {
	return redisClient
}
