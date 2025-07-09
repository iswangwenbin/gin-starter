package service

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseService struct {
	DB     *gorm.DB
	Cache  *redis.Client
	Logger *zap.Logger
	Ctx    context.Context
}

func NewBaseService(db *gorm.DB, cache *redis.Client, logger *zap.Logger) *BaseService {
	return &BaseService{
		DB:     db,
		Cache:  cache,
		Logger: logger,
		Ctx:    context.Background(),
	}
}

func (bs *BaseService) WithContext(ctx context.Context) *BaseService {
	return &BaseService{
		DB:     bs.DB,
		Cache:  bs.Cache,
		Logger: bs.Logger,
		Ctx:    ctx,
	}
}