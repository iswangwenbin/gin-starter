package worker

import (
	"context"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/iswangwenbin/gin-starter/internal/repository"
	"github.com/iswangwenbin/gin-starter/internal/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// InstallEventWorker 安装事件处理工作者
type InstallEventWorker struct {
	redis     *redis.Client
	clickHouse clickhouse.Conn
	logger    *zap.Logger
	consumer  *service.InstallEventConsumer
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewInstallEventWorker(redis *redis.Client, clickHouse clickhouse.Conn, logger *zap.Logger) *InstallEventWorker {
	ctx, cancel := context.WithCancel(context.Background())
	
	// 创建 Repository
	installEventRepo := repository.NewInstallEventRepository(clickHouse)
	
	// 创建 Consumer
	consumer := service.NewInstallEventConsumer(redis, installEventRepo, logger)
	
	return &InstallEventWorker{
		redis:     redis,
		clickHouse: clickHouse,
		logger:    logger,
		consumer:  consumer,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start 启动工作者
func (w *InstallEventWorker) Start() error {
	w.logger.Info("Starting install event worker...")
	
	// 启动消费者
	if err := w.consumer.Start(); err != nil {
		w.logger.Error("Failed to start consumer", zap.Error(err))
		return err
	}
	
	w.logger.Info("Install event worker started successfully")
	return nil
}

// Stop 停止工作者
func (w *InstallEventWorker) Stop() {
	w.logger.Info("Stopping install event worker...")
	
	// 停止消费者
	w.consumer.Stop()
	
	// 取消上下文
	w.cancel()
	
	// 等待所有 goroutine 完成
	w.wg.Wait()
	
	w.logger.Info("Install event worker stopped")
}

// GetStatus 获取工作者状态
func (w *InstallEventWorker) GetStatus() (map[string]interface{}, error) {
	pendingCount, err := w.consumer.GetPendingCount()
	if err != nil {
		return nil, err
	}
	
	streamLength, err := w.consumer.GetStreamLength()
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"pending_count": pendingCount,
		"stream_length": streamLength,
		"consumer_group": service.InstallEventConsumerGroup,
		"consumer_name": service.InstallEventConsumerName,
		"stream_key": service.InstallEventStreamKey,
	}, nil
}