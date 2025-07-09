package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	Redis  *redis.Client
	Limit  int
	Window time.Duration
}

func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		Redis:  redisClient,
		Limit:  limit,
		Window: window,
	}
}

func (rl *RateLimiter) Allow(key string) (bool, error) {
	ctx := context.Background()
	now := time.Now()
	windowStart := now.Truncate(rl.Window)
	
	pipe := rl.Redis.Pipeline()
	
	// 使用滑动窗口计数器
	countKey := fmt.Sprintf("rate_limit:%s:%d", key, windowStart.Unix())
	
	// 增加计数
	pipe.Incr(ctx, countKey)
	// 设置过期时间
	pipe.Expire(ctx, countKey, rl.Window)
	
	results, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	
	count := results[0].(*redis.IntCmd).Val()
	return count <= int64(rl.Limit), nil
}

func (rl *RateLimiter) GetCurrentCount(key string) (int64, error) {
	ctx := context.Background()
	now := time.Now()
	windowStart := now.Truncate(rl.Window)
	countKey := fmt.Sprintf("rate_limit:%s:%d", key, windowStart.Unix())
	
	return rl.Redis.Get(ctx, countKey).Int64()
}

func RateLimit(rateLimiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用客户端IP作为限制key
		key := c.ClientIP()
		
		allowed, err := rateLimiter.Allow(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "rate limit check failed",
				"data":    nil,
			})
			c.Abort()
			return
		}
		
		if !allowed {
			c.Header("X-RateLimit-Limit", strconv.Itoa(rateLimiter.Limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rateLimiter.Window).Unix(), 10))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "rate limit exceeded",
				"data":    nil,
			})
			c.Abort()
			return
		}
		
		// 设置响应头
		count, _ := rateLimiter.GetCurrentCount(key)
		remaining := rateLimiter.Limit - int(count)
		if remaining < 0 {
			remaining = 0
		}
		
		c.Header("X-RateLimit-Limit", strconv.Itoa(rateLimiter.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rateLimiter.Window).Unix(), 10))
		
		c.Next()
	}
}

func RateLimitByUser(rateLimiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			// 如果没有用户ID，使用IP限流
			key := c.ClientIP()
			allowed, err := rateLimiter.Allow(key)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "rate limit check failed",
					"data":    nil,
				})
				c.Abort()
				return
			}
			
			if !allowed {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"code":    429,
					"message": "rate limit exceeded",
					"data":    nil,
				})
				c.Abort()
				return
			}
		} else {
			// 使用用户ID限流
			key := fmt.Sprintf("user:%v", userID)
			allowed, err := rateLimiter.Allow(key)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "rate limit check failed",
					"data":    nil,
				})
				c.Abort()
				return
			}
			
			if !allowed {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"code":    429,
					"message": "rate limit exceeded",
					"data":    nil,
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}