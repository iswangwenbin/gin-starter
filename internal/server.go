package internal

import (
	"context"
	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iswangwenbin/gin-starter/internal/middleware"
	"github.com/iswangwenbin/gin-starter/pkg/configx"
	"github.com/iswangwenbin/gin-starter/pkg/databasex"
	"github.com/iswangwenbin/gin-starter/pkg/redisx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	Engine      *gin.Engine   // 框架
	DB          *gorm.DB      // 数据库
	Cache       *redis.Client // Redis
	Environment string        // 运行环境
	logger      *zap.Logger   // 日志

	startPProf    bool // 是否初始化PProf
	startDatabase bool // 是否初始化数据库
	startDebug    bool // 是否初始化调试模式
	startCache    bool // 是否初始化Redis
}

func NewServer(env string, options ...Option) (*Server, error) {
	s := &Server{}
	s.Environment = env

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, "option error")
		}
	}

	// Engine
	if s.Engine == nil {
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
			return nil, err
		}

		// 日志初始化
		if env == "production" {
			s.logger, _ = zap.NewProduction()
		} else {
			s.logger, _ = zap.NewDevelopment()
		}

		// 设置全局logger，这样zap.S()才能正常工作
		zap.ReplaceGlobals(s.logger)

		// 添加中间件
		s.Engine.Use(middleware.RequestID())
		s.Engine.Use(middleware.Security())
		s.Engine.Use(middleware.HidePoweredBy())
		s.Engine.Use(middleware.CORS())
		s.Engine.Use(ginzap.Ginzap(s.logger, time.RFC3339, true))
		s.Engine.Use(middleware.ErrorHandler(s.logger))
		
		// 设置404和405处理
		s.Engine.NoRoute(middleware.NotFoundHandler())
		s.Engine.NoMethod(middleware.MethodNotAllowedHandler())
	}

	// Database
	if s.startDatabase {
		s.DB = databasex.NewDB()
		s.logger.Info("Database Enable")
	}

	// Cache
	if s.startCache {
		s.Cache = redisx.GetRedis()
		s.logger.Info("Redis Cache Enable")
	}
	return s, nil
}

func (s *Server) Start() {
	// 路由
	s.routes()
	// 监听端口
	s.listen()
}

// 监听端口
func (s *Server) listen() {
	cfg := configx.GetConfig()
	if cfg == nil {
		log.Fatal("Config not loaded")
	}
	
	srv := &http.Server{
		Addr:           cfg.GetServerAddress(),
		Handler:        s.Engine,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}
	// 在一个 goroutine 中启动服务器，这样它就不会阻塞下面的优雅关机处理
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	s.logger.Info("Server exiting")
}
