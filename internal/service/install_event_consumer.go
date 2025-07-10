package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iswangwenbin/gin-starter/internal/model"
	"github.com/iswangwenbin/gin-starter/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	InstallEventConsumerGroup = "install_events_consumer_group"
	InstallEventConsumerName  = "install_events_consumer"
	BatchSize                 = 100 // 批量处理大小
	BatchTimeout              = 5 * time.Second
)

type InstallEventConsumer struct {
	redis            *redis.Client
	installEventRepo repository.InstallEventRepository
	logger           *zap.Logger
	ctx              context.Context
	cancel           context.CancelFunc
}

func NewInstallEventConsumer(redis *redis.Client, installEventRepo repository.InstallEventRepository, logger *zap.Logger) *InstallEventConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &InstallEventConsumer{
		redis:            redis,
		installEventRepo: installEventRepo,
		logger:           logger,
		ctx:              ctx,
		cancel:           cancel,
	}
}

// Start 启动消费者
func (c *InstallEventConsumer) Start() error {
	// 创建消费者组（如果不存在）
	err := c.redis.XGroupCreateMkStream(c.ctx, InstallEventStreamKey, InstallEventConsumerGroup, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		c.logger.Error("Failed to create consumer group", zap.Error(err))
		return err
	}

	c.logger.Info("Install event consumer started",
		zap.String("stream", InstallEventStreamKey),
		zap.String("group", InstallEventConsumerGroup),
		zap.String("consumer", InstallEventConsumerName))

	// 启动消费循环
	go c.consumeLoop()

	return nil
}

// Stop 停止消费者
func (c *InstallEventConsumer) Stop() {
	c.logger.Info("Stopping install event consumer...")
	c.cancel()
}

// consumeLoop 消费循环
func (c *InstallEventConsumer) consumeLoop() {
	batch := make([]*model.InstallEvent, 0, BatchSize)
	messageIDs := make([]string, 0, BatchSize)
	
	ticker := time.NewTicker(BatchTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			// 处理剩余的批次
			if len(batch) > 0 {
				c.processBatch(batch, messageIDs)
			}
			c.logger.Info("Install event consumer stopped")
			return

		case <-ticker.C:
			// 定时处理批次
			if len(batch) > 0 {
				c.processBatch(batch, messageIDs)
				batch = batch[:0]
				messageIDs = messageIDs[:0]
			}

		default:
			// 读取消息
			messages, err := c.redis.XReadGroup(c.ctx, &redis.XReadGroupArgs{
				Group:    InstallEventConsumerGroup,
				Consumer: InstallEventConsumerName,
				Streams:  []string{InstallEventStreamKey, ">"},
				Count:    10,
				Block:    time.Second,
			}).Result()

			if err != nil {
				if err != redis.Nil && err != context.DeadlineExceeded {
					c.logger.Error("Failed to read from stream", zap.Error(err))
				}
				continue
			}

			// 处理消息
			for _, stream := range messages {
				for _, message := range stream.Messages {
					event, err := c.parseMessage(message)
					if err != nil {
						c.logger.Error("Failed to parse message",
							zap.String("message_id", message.ID),
							zap.Error(err))
						// 确认消息以避免重复处理
						c.ackMessage(message.ID)
						continue
					}

					batch = append(batch, event)
					messageIDs = append(messageIDs, message.ID)

					// 批次满了，立即处理
					if len(batch) >= BatchSize {
						c.processBatch(batch, messageIDs)
						batch = batch[:0]
						messageIDs = messageIDs[:0]
						ticker.Reset(BatchTimeout) // 重置定时器
					}
				}
			}
		}
	}
}

// parseMessage 解析消息
func (c *InstallEventConsumer) parseMessage(message redis.XMessage) (*model.InstallEvent, error) {
	eventDataStr, ok := message.Values["event_data"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid event_data format")
	}

	var req model.CreateInstallEventRequest
	if err := json.Unmarshal([]byte(eventDataStr), &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	// 转换为 InstallEvent
	event := &model.InstallEvent{
		AppID:            req.AppID,
		AppName:          req.AppName,
		AppVersion:       req.AppVersion,
		AppType:          req.AppType,
		EventID:          req.EventID,
		EventDate:        req.EventTime.Truncate(24 * time.Hour),
		EventTime:        req.EventTime,
		DeviceID:         req.DeviceID,
		ChannelID:        req.ChannelID,
		InstallIP:        req.InstallIP,
		InstallType:      req.InstallType,
		InstallResult:    req.InstallResult,
		OSLanguage:       req.OSLanguage,
		OSTimezone:       req.OSTimezone,
		OSName:           req.OSName,
		OSVersion:        req.OSVersion,
		OSBuild:          req.OSBuild,
		OSFamily:         req.OSFamily,
		SignatureStatus:  req.SignatureStatus,
		SignatureVersion: req.SignatureVersion,
		SignatureParams:  req.SignatureParams,
	}

	return event, nil
}

// processBatch 批量处理事件
func (c *InstallEventConsumer) processBatch(batch []*model.InstallEvent, messageIDs []string) {
	if len(batch) == 0 {
		return
	}

	c.logger.Info("Processing install events batch", zap.Int("count", len(batch)))

	// 批量写入 ClickHouse
	if err := c.installEventRepo.CreateBatch(c.ctx, batch); err != nil {
		c.logger.Error("Failed to write batch to ClickHouse",
			zap.Int("count", len(batch)),
			zap.Error(err))
		return
	}

	// 确认所有消息
	for _, messageID := range messageIDs {
		c.ackMessage(messageID)
	}

	c.logger.Info("Install events batch processed successfully", zap.Int("count", len(batch)))
}

// ackMessage 确认消息
func (c *InstallEventConsumer) ackMessage(messageID string) {
	if err := c.redis.XAck(c.ctx, InstallEventStreamKey, InstallEventConsumerGroup, messageID).Err(); err != nil {
		c.logger.Error("Failed to ack message",
			zap.String("message_id", messageID),
			zap.Error(err))
	}
}

// GetPendingCount 获取待处理消息数量
func (c *InstallEventConsumer) GetPendingCount() (int64, error) {
	info, err := c.redis.XPending(c.ctx, InstallEventStreamKey, InstallEventConsumerGroup).Result()
	if err != nil {
		return 0, err
	}
	return info.Count, nil
}

// GetStreamLength 获取流长度
func (c *InstallEventConsumer) GetStreamLength() (int64, error) {
	return c.redis.XLen(c.ctx, InstallEventStreamKey).Result()
}