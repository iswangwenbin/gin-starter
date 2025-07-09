package redisx

import (
	"context"
	"log"
	"time"

	"github.com/iswangwenbin/gin-starter/pkg/configx"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func NewRedis() *redis.Client {
	cfg := configx.GetConfig()
	if cfg == nil {
		log.Fatal("Config not loaded")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddress(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		DialTimeout:  10 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	redisClient = rdb
	return rdb
}

func GetRedis() *redis.Client {
	if redisClient == nil {
		redisClient = NewRedis()
	}
	return redisClient
}

func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

func GetStats() *redis.PoolStats {
	if redisClient != nil {
		return redisClient.PoolStats()
	}
	return nil
}
