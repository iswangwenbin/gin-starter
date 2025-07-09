package server

import (
	"context"

	"github.com/iswangwenbin/gin-starter/internal/grpc/protobuf"
)

// HealthServer 健康检查服务的 gRPC 服务端实现
type HealthServer struct {
	protobuf.UnimplementedHealthServiceServer
}

// NewHealthServer 创建健康检查服务的 gRPC 服务端
func NewHealthServer() *HealthServer {
	return &HealthServer{}
}

// Check 健康检查
func (s *HealthServer) Check(ctx context.Context, req *protobuf.HealthCheckRequest) (*protobuf.HealthCheckResponse, error) {
	// 这里可以添加实际的健康检查逻辑
	// 例如：检查数据库连接、Redis连接等
	
	return &protobuf.HealthCheckResponse{
		Status:  protobuf.HealthCheckResponse_SERVING,
		Message: "Service is healthy",
	}, nil
}