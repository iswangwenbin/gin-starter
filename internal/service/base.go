package service

import (
	"context"

	"github.com/iswangwenbin/gin-starter/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type BaseService struct {
	Repo   *repository.RepositoryManager
	Cache  *redis.Client
	Logger *zap.Logger
	Ctx    context.Context
}

func NewBaseService(repo *repository.RepositoryManager, cache *redis.Client, logger *zap.Logger) *BaseService {
	return &BaseService{
		Repo:   repo,
		Cache:  cache,
		Logger: logger,
		Ctx:    context.Background(),
	}
}

func (bs *BaseService) WithContext(ctx context.Context) *BaseService {
	return &BaseService{
		Repo:   bs.Repo,
		Cache:  bs.Cache,
		Logger: bs.Logger,
		Ctx:    ctx,
	}
}