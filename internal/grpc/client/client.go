package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/iswangwenbin/gin-starter/internal/grpc/protobuf"
	"go.uber.org/zap"
)

// ClientManager gRPC 客户端管理器
type ClientManager struct {
	address    string
	conn       *grpc.ClientConn
	logger     *zap.Logger
	userClient *UserClient
	healthClient protobuf.HealthServiceClient
}

// NewClientManager 创建 gRPC 客户端管理器
func NewClientManager(address string, logger *zap.Logger) (*ClientManager, error) {
	// 连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	// 创建连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &ClientManager{
		address:      address,
		conn:         conn,
		logger:       logger,
		healthClient: protobuf.NewHealthServiceClient(conn),
	}, nil
}

// Close 关闭客户端连接
func (m *ClientManager) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

// UserClient 获取用户客户端
func (m *ClientManager) UserClient() *UserClient {
	if m.userClient == nil {
		m.userClient = &UserClient{
			client: protobuf.NewUserServiceClient(m.conn),
			conn:   m.conn,
			logger: m.logger,
		}
	}
	return m.userClient
}

// HealthClient 获取健康检查客户端
func (m *ClientManager) HealthClient() protobuf.HealthServiceClient {
	return m.healthClient
}

// HealthCheck 健康检查
func (m *ClientManager) HealthCheck(ctx context.Context, service string) (*protobuf.HealthCheckResponse, error) {
	req := &protobuf.HealthCheckRequest{
		Service: service,
	}
	return m.healthClient.Check(ctx, req)
}

// IsHealthy 检查服务是否健康
func (m *ClientManager) IsHealthy(ctx context.Context, service string) bool {
	resp, err := m.HealthCheck(ctx, service)
	if err != nil {
		m.logger.Error("Health check failed", zap.Error(err))
		return false
	}
	return resp.Status == protobuf.HealthCheckResponse_SERVING
}