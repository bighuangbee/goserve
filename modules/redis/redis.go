package redis

import (
	"github.com/go-redis/redis"
	"goserve/modules/config"
	"goserve/modules/loger"
)

var Redis *redis.Client

func Setup(){

	Redis = redis.NewClient(&redis.Options{
		Addr:     config.ConfigData.RedisAddr,
		Password: config.ConfigData.RedisPassword,
		DB:       config.ConfigData.RedisDefaultDB,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		loger.Error("Redis Client SetUp Failed ### ", err)
		return
	}

	loger.Info("Redis Client SetUp Success...")
}