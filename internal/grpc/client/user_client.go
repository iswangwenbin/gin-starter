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

// UserClient 用户服务的 gRPC 客户端
type UserClient struct {
	client protobuf.UserServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

// NewUserClient 创建用户服务的 gRPC 客户端
func NewUserClient(address string, logger *zap.Logger) (*UserClient, error) {
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

	// 创建客户端
	client := protobuf.NewUserServiceClient(conn)

	return &UserClient{
		client: client,
		conn:   conn,
		logger: logger,
	}, nil
}

// Close 关闭客户端连接
func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateUser 创建用户
func (c *UserClient) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.CreateUserResponse, error) {
	return c.client.CreateUser(ctx, req)
}

// GetUser 获取用户
func (c *UserClient) GetUser(ctx context.Context, req *protobuf.GetUserRequest) (*protobuf.GetUserResponse, error) {
	return c.client.GetUser(ctx, req)
}

// UpdateUser 更新用户
func (c *UserClient) UpdateUser(ctx context.Context, req *protobuf.UpdateUserRequest) (*protobuf.UpdateUserResponse, error) {
	return c.client.UpdateUser(ctx, req)
}

// DeleteUser 删除用户
func (c *UserClient) DeleteUser(ctx context.Context, req *protobuf.DeleteUserRequest) (*protobuf.DeleteUserResponse, error) {
	return c.client.DeleteUser(ctx, req)
}

// ListUsers 获取用户列表
func (c *UserClient) ListUsers(ctx context.Context, req *protobuf.ListUsersRequest) (*protobuf.ListUsersResponse, error) {
	return c.client.ListUsers(ctx, req)
}

// Login 用户登录
func (c *UserClient) Login(ctx context.Context, req *protobuf.LoginRequest) (*protobuf.LoginResponse, error) {
	return c.client.Login(ctx, req)
}

// ChangePassword 修改密码
func (c *UserClient) ChangePassword(ctx context.Context, req *protobuf.ChangePasswordRequest) (*protobuf.ChangePasswordResponse, error) {
	return c.client.ChangePassword(ctx, req)
}