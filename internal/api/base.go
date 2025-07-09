package api

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseController struct {
	DB     *gorm.DB
	Cache  *redis.Client
	Logger *zap.Logger
}

func NewBaseController(db *gorm.DB, cache *redis.Client, logger *zap.Logger) *BaseController {
	return &BaseController{
		DB:     db,
		Cache:  cache,
		Logger: logger,
	}
}

func (bc *BaseController) GetLogger(c *gin.Context) *zap.Logger {
	return bc.Logger.With(
		zap.String("request_id", c.GetString("request_id")),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)
}