package server

import (
	"context"

	"github.com/iswangwenbin/gin-starter/internal/grpc/protobuf"
	"github.com/iswangwenbin/gin-starter/internal/model"
	"github.com/iswangwenbin/gin-starter/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InstallEventServer struct {
	protobuf.UnimplementedInstallEventServiceServer
	installEventService *service.InstallEventService
	logger              *zap.Logger
}

func NewInstallEventServer(installEventService *service.InstallEventService, logger *zap.Logger) *InstallEventServer {
	return &InstallEventServer{
		installEventService: installEventService,
		logger:              logger,
	}
}

// 创建单个安装事件
func (s *InstallEventServer) CreateInstallEvent(ctx context.Context, req *protobuf.CreateInstallEventRequest) (*protobuf.CreateInstallEventResponse, error) {
	// 转换 protobuf 请求到内部模型
	createReq := &model.CreateInstallEventRequest{
		AppID:            req.AppId,
		AppName:          req.AppName,
		AppVersion:       req.AppVersion,
		AppType:          model.AppType(req.AppType),
		EventID:          req.EventId,
		EventTime:        req.EventTime.AsTime(),
		DeviceID:         req.DeviceId,
		ChannelID:        req.ChannelId,
		InstallIP:        req.InstallIp,
		InstallType:      model.InstallType(req.InstallType),
		InstallResult:    model.InstallResult(req.InstallResult),
		OSLanguage:       req.OsLanguage,
		OSTimezone:       req.OsTimezone,
		OSName:           req.OsName,
		OSVersion:        req.OsVersion,
		OSBuild:          req.OsBuild,
		OSFamily:         req.OsFamily,
		SignatureStatus:  uint8(req.SignatureStatus),
		SignatureVersion: req.SignatureVersion,
		SignatureParams:  req.SignatureParams,
	}

	// 调用服务层
	if err := s.installEventService.Create(ctx, createReq); err != nil {
		s.logger.Error("Failed to create install event via gRPC",
			zap.String("event_id", req.EventId),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to create install event: %v", err)
	}

	return &protobuf.CreateInstallEventResponse{
		Success: true,
		Message: "Install event created successfully",
	}, nil
}

// 批量创建安装事件
func (s *InstallEventServer) CreateInstallEventBatch(ctx context.Context, req *protobuf.CreateInstallEventBatchRequest) (*protobuf.CreateInstallEventBatchResponse, error) {
	if len(req.Events) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "No events provided")
	}

	// 转换 protobuf 请求到内部模型
	createReqs := make([]*model.CreateInstallEventRequest, 0, len(req.Events))
	for _, event := range req.Events {
		createReq := &model.CreateInstallEventRequest{
			AppID:            event.AppId,
			AppName:          event.AppName,
			AppVersion:       event.AppVersion,
			AppType:          model.AppType(event.AppType),
			EventID:          event.EventId,
			EventTime:        event.EventTime.AsTime(),
			DeviceID:         event.DeviceId,
			ChannelID:        event.ChannelId,
			InstallIP:        event.InstallIp,
			InstallType:      model.InstallType(event.InstallType),
			InstallResult:    model.InstallResult(event.InstallResult),
			OSLanguage:       event.OsLanguage,
			OSTimezone:       event.OsTimezone,
			OSName:           event.OsName,
			OSVersion:        event.OsVersion,
			OSBuild:          event.OsBuild,
			OSFamily:         event.OsFamily,
			SignatureStatus:  uint8(event.SignatureStatus),
			SignatureVersion: event.SignatureVersion,
			SignatureParams:  event.SignatureParams,
		}
		createReqs = append(createReqs, createReq)
	}

	// 调用服务层批量创建
	if err := s.installEventService.CreateBatch(ctx, createReqs); err != nil {
		s.logger.Error("Failed to create install events batch via gRPC",
			zap.Int("count", len(req.Events)),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to create install events batch: %v", err)
	}

	return &protobuf.CreateInstallEventBatchResponse{
		Success:     true,
		Message:     "Install events batch created successfully",
		ProcessedCount: int32(len(req.Events)),
	}, nil
}