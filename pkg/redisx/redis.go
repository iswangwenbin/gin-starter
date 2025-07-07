package redisx

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
)

var redisConnection *redis.Client

func NewRedis() *redis.Client {
	ctx := context.Background()
	addr := viper.GetString("Redis.addr")
	password := viper.GetString("Redis.password")
	redisConnection := redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		})
	_, err := redisConnection.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("NewRedis.error:%v", err)
		return nil
	}
	return redisConnection
}

func GetRedis() *redis.Client {
	if redisConnection == nil {
		redisConnection = NewRedis()
	}
	return redisConnection
}
