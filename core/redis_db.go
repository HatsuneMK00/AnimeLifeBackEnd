package core

import (
	"AnimeLifeBackEnd/global"
	"github.com/redis/go-redis/v9"
)

func InitRedisDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Addr,
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.DB,
	})
	return rdb
}
