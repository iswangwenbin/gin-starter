package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/iswangwenbin/gin-starter/internal/grpc/protobuf"
	"github.com/iswangwenbin/gin-starter/internal/repository"
	"github.com/iswangwenbin/gin-starter/internal/service"
	"github.com/iswangwenbin/gin-starter/pkg/configx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server gRPC 服务器
type Server struct {
	grpcServer *grpc.Server
	config     *configx.Config
	logger     *zap.Logger
	db         *gorm.DB
	cache      *redis.Client
}

// NewServer 创建 gRPC 服务器
func NewServer(config *configx.Config, logger *zap.Logger, db *gorm.DB, cache *redis.Client) *Server {
	return &Server{
		config: config,
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

// Start 启动 gRPC 服务器
func (s *Server) Start(ctx context.Context) error {
	// 创建 gRPC 服务器选项
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second,
			MaxConnectionAge:      30 * time.Second,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  5 * time.Second,
			Timeout:               1 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	// 创建 gRPC 服务器
	s.grpcServer = grpc.NewServer(opts...)

	// 注册服务
	s.registerServices()

	// 启用反射（开发环境）
	if s.config.Debug {
		reflection.Register(s.grpcServer)
	}

	// 监听端口
	addr := s.config.GetGRPCAddress()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.logger.Info("gRPC server starting", zap.String("address", addr))

	// 启动服务器
	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			s.logger.Error("gRPC server failed", zap.Error(err))
		}
	}()

	// 等待上下文取消
	<-ctx.Done()

	// 优雅关闭
	s.logger.Info("gRPC server shutting down...")
	s.grpcServer.GracefulStop()

	return nil
}

// Stop 停止 gRPC 服务器
func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// registerServices 注册 gRPC 服务
func (s *Server) registerServices() {
	// 创建服务层依赖
	repo := repository.NewRepository(s.db)
	baseService := service.NewBaseService(repo, s.cache, s.logger)

	// 创建用户服务
	userService := service.NewUserService(baseService)
	userServer := NewUserServer(userService)

	// 创建健康检查服务
	healthServer := NewHealthServer()

	// 注册服务
	protobuf.RegisterUserServiceServer(s.grpcServer, userServer)
	protobuf.RegisterHealthServiceServer(s.grpcServer, healthServer)

	s.logger.Info("gRPC services registered")
}