package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iswangwenbin/gin-starter/internal/model"
	"github.com/iswangwenbin/gin-starter/pkg/errorsx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	InstallEventStreamKey = "install_events_stream"
)

type InstallEventService struct {
	redis  *redis.Client
	logger *zap.Logger
}

func NewInstallEventService(redis *redis.Client, logger *zap.Logger) *InstallEventService {
	return &InstallEventService{
		redis:  redis,
		logger: logger,
	}
}

// Create 创建单个安装事件 - 写入 Redis Stream
func (s *InstallEventService) Create(ctx context.Context, req *model.CreateInstallEventRequest) error {

	fmt.Printf("create install event request: %+v\n", req)

	// 数据验证
	if req.EventID == "" {
		return errorsx.New(errorsx.CodeBadRequest, "EventID is required")
	}

	// 序列化请求数据
	eventData, err := json.Marshal(req)
	if err != nil {
		s.logger.Error("Failed to marshal install event",
			zap.String("event_id", req.EventID),
			zap.Error(err))
		return errorsx.NewWithError(errorsx.CodeInternalServerError, "Failed to serialize event data", err)
	}

	// 写入 Redis Stream
	result := s.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: InstallEventStreamKey,
		MaxLen: 100000, // 保留最近10万条记录
		Approx: true,   // 使用近似裁剪，性能更好
		Values: map[string]interface{}{
			"event_id":   req.EventID,
			"app_id":     req.AppID,
			"device_id":  req.DeviceID,
			"event_data": string(eventData),
			"created_at": time.Now().Unix(),
		},
	})

	if result.Err() != nil {
		s.logger.Error("Failed to add install event to stream",
			zap.String("event_id", req.EventID),
			zap.String("app_id", req.AppID),
			zap.Error(result.Err()))
		return errorsx.NewWithError(errorsx.CodeInternalServerError, "Failed to queue event", result.Err())
	}

	s.logger.Info("Install event queued successfully",
		zap.String("event_id", req.EventID),
		zap.String("app_id", req.AppID),
		zap.String("device_id", req.DeviceID),
		zap.String("stream_id", result.Val()))

	return nil
}

// CreateBatch 批量创建安装事件 - 写入 Redis Stream
func (s *InstallEventService) CreateBatch(ctx context.Context, requests []*model.CreateInstallEventRequest) error {
	if len(requests) == 0 {
		return nil
	}

	// 使用 Pipeline 批量写入
	pipe := s.redis.Pipeline()
	validCount := 0
	createdAt := time.Now().Unix()

	for _, req := range requests {
		// 数据验证
		if req.EventID == "" {
			s.logger.Warn("Skipping event with empty EventID", zap.String("app_id", req.AppID))
			continue
		}

		// 序列化请求数据
		eventData, err := json.Marshal(req)
		if err != nil {
			s.logger.Warn("Skipping event due to marshal error",
				zap.String("event_id", req.EventID),
				zap.Error(err))
			continue
		}

		// 添加到 Pipeline
		pipe.XAdd(ctx, &redis.XAddArgs{
			Stream: InstallEventStreamKey,
			MaxLen: 100000,
			Approx: true,
			Values: map[string]interface{}{
				"event_id":   req.EventID,
				"app_id":     req.AppID,
				"device_id":  req.DeviceID,
				"event_data": string(eventData),
				"created_at": createdAt,
			},
		})
		validCount++
	}

	if validCount == 0 {
		return errorsx.New(errorsx.CodeBadRequest, "No valid events to create")
	}

	// 执行 Pipeline
	results, err := pipe.Exec(ctx)
	if err != nil {
		s.logger.Error("Failed to execute batch pipeline",
			zap.Int("count", validCount),
			zap.Error(err))
		return errorsx.NewWithError(errorsx.CodeInternalServerError, "Failed to queue events batch", err)
	}

	// 检查结果
	successCount := 0
	for _, result := range results {
		if result.Err() == nil {
			successCount++
		}
	}

	s.logger.Info("Install events batch queued",
		zap.Int("total", len(requests)),
		zap.Int("valid", validCount),
		zap.Int("success", successCount))

	return nil
}
