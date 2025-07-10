package core

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/iswangwenbin/gin-starter/internal/grpc/server"
	"github.com/iswangwenbin/gin-starter/internal/middleware"
	"github.com/iswangwenbin/gin-starter/pkg/clickhousex"
	"github.com/iswangwenbin/gin-starter/pkg/configx"
	"github.com/iswangwenbin/gin-starter/pkg/databasex"
	"github.com/iswangwenbin/gin-starter/pkg/redisx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	Engine           *gin.Engine                   // HTTP 框架
	GRPCServer       *server.Server                // gRPC 服务器
	DB               *gorm.DB                      // 数据库
	Cache            *redis.Client                 // Redis
	ClickHouse       clickhouse.Conn               // ClickHouse
	Environment      string                        // 运行环境
	logger           *zap.Logger                   // 日志

	startPProf      bool // 是否初始化PProf
	startDatabase   bool // 是否初始化数据库
	startDebug      bool // 是否初始化调试模式
	startCache      bool // 是否初始化Redis
	startGRPC       bool // 是否启动gRPC服务器
	startClickHouse bool // 是否初始化ClickHouse
}

func NewServer(env string, options ...Option) (*Server, error) {
	s := &Server{Environment: env}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, "option error")
		}
	}

	if err := s.initEngine(env); err != nil {
		return nil, err
	}

	s.initComponents()

	return s, nil
}

func (s *Server) initEngine(env string) error {
	if s.Engine != nil {
		return nil
	}

	switch env {
	case "production":
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	s.Engine = gin.New()
	err := s.Engine.SetTrustedProxies([]string{"127.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12", "10.0.0.0/8"})
	if err != nil {
		return err
	}

	s.initLogger(env)
	s.setupMiddleware()

	return nil
}

func (s *Server) initLogger(env string) {
	if env == "production" {
		s.logger, _ = zap.NewProduction()
	} else {
		s.logger, _ = zap.NewDevelopment()
	}
	zap.ReplaceGlobals(s.logger)
}

func (s *Server) setupMiddleware() {
	s.Engine.Use(middleware.RequestID())
	s.Engine.Use(middleware.Security())
	s.Engine.Use(middleware.HidePoweredBy())
	s.Engine.Use(middleware.CORS())
	s.Engine.Use(ginzap.Ginzap(s.logger, time.RFC3339, true))
	s.Engine.Use(middleware.ErrorHandler(s.logger))

	s.Engine.NoRoute(middleware.NotFoundHandler())
	s.Engine.NoMethod(middleware.MethodNotAllowedHandler())
}

func (s *Server) initComponents() {
	if s.startDatabase {
		s.DB = databasex.NewDB()
		s.logger.Info("Database Enable")
	}

	if s.startCache {
		s.Cache = redisx.GetRedis()
		s.logger.Info("Redis Cache Enable")
	}

	if s.startClickHouse {
		s.ClickHouse = clickhousex.NewClickHouse()
		s.logger.Info("ClickHouse Enable")
	}

	if s.startGRPC {
		cfg := configx.GetConfig()
		if cfg != nil && cfg.GRPC.Enabled {
			s.GRPCServer = server.NewServer(cfg, s.logger, s.DB, s.Cache)
			s.logger.Info("gRPC Server Enable")
		}
	}

}

// Logger 获取日志实例
func (s *Server) Logger() *zap.Logger {
	return s.logger
}

// Start 启动服务器
func (s *Server) Start() error {

	// 路由
	s.routes()
	
	// 监听端口（阻塞）
	s.listen()
	
	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	
	s.logger.Info("Server stopped")
	return nil
}

// 监听端口
func (s *Server) listen() {
	cfg := configx.GetConfig()
	if cfg == nil {
		log.Fatal("Config not loaded")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	httpSrv := s.startHTTPServer(cfg, &wg)
	s.startGRPCServer(ctx, &wg)

	s.logServerStatus(cfg)
	s.handleShutdown(cancel, httpSrv, &wg)
}

func (s *Server) startHTTPServer(cfg *configx.Config, wg *sync.WaitGroup) *http.Server {
	httpSrv := &http.Server{
		Addr:           cfg.GetServerAddress(),
		Handler:        s.Engine,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.Info("HTTP server starting", zap.String("address", httpSrv.Addr))

		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("HTTP server error", zap.Error(err))
		}
	}()

	return httpSrv
}

func (s *Server) startGRPCServer(ctx context.Context, wg *sync.WaitGroup) {
	if s.GRPCServer == nil {
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.GRPCServer.Start(ctx); err != nil {
			s.logger.Error("gRPC server error", zap.Error(err))
		}
	}()
}

func (s *Server) logServerStatus(cfg *configx.Config) {
	s.logger.Info("Servers started successfully",
		zap.String("http_address", cfg.GetServerAddress()),
		zap.String("grpc_address", cfg.GetGRPCAddress()),
		zap.Bool("grpc_enabled", cfg.GRPC.Enabled && s.GRPCServer != nil),
	)
}

func (s *Server) handleShutdown(cancel context.CancelFunc, httpSrv *http.Server, wg *sync.WaitGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.logger.Info("Received shutdown signal, gracefully shutting down...")


	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("HTTP server shutdown error", zap.Error(err))
	} else {
		s.logger.Info("HTTP server shut down successfully")
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("All servers shut down successfully")
	case <-time.After(30 * time.Second):
		s.logger.Warn("Shutdown timeout, forcing exit")
	}
}