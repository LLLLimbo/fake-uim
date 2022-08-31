package settings

import (
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client = nil

func RedisCli() *redis.Client {
	return rdb
}

func InitRdb() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.2.102:6379",
		Password: "seeiner",
		DB:       0,
	})

}
