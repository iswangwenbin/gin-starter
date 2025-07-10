package core

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// Lifecycle 应用生命周期管理
type Lifecycle struct {
	server *Server
	logger *zap.Logger
	done   chan struct{}
}

func NewLifecycle(server *Server) *Lifecycle {
	return &Lifecycle{
		server: server,
		logger: server.Logger(),
		done:   make(chan struct{}),
	}
}

// Run 运行应用程序（阻塞直到收到停止信号）
func (l *Lifecycle) Run() error {
	l.logger.Info("Application starting...")

	// 启动服务器（非阻塞）
	go func() {
		if err := l.server.Start(); err != nil {
			l.logger.Error("Server start failed", zap.Error(err))
			l.Stop()
		}
	}()

	// 等待停止信号
	l.waitForShutdown()

	// 优雅关闭
	l.logger.Info("Application shutting down...")
	if err := l.server.Stop(); err != nil {
		l.logger.Error("Server stop failed", zap.Error(err))
		return err
	}

	l.logger.Info("Application stopped successfully")
	return nil
}

// Stop 手动停止应用程序
func (l *Lifecycle) Stop() {
	select {
	case <-l.done:
		// 已经停止
	default:
		close(l.done)
	}
}

// waitForShutdown 等待关闭信号
func (l *Lifecycle) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		l.logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	case <-l.done:
		l.logger.Info("Received stop request")
	}
}