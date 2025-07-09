package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iswangwenbin/gin-starter/internal/grpc/protobuf"
	"github.com/iswangwenbin/gin-starter/internal/middleware"
	"github.com/iswangwenbin/gin-starter/internal/model"
	"github.com/iswangwenbin/gin-starter/internal/service"
	"github.com/iswangwenbin/gin-starter/pkg/errorsx"
)

// UserServer 用户服务的 gRPC 服务端实现
type UserServer struct {
	protobuf.UnimplementedUserServiceServer
	userService *service.UserService
}

// NewUserServer 创建用户服务的 gRPC 服务端
func NewUserServer(userService *service.UserService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.CreateUserResponse, error) {
	// 转换请求
	createReq := &model.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
	}

	// 设置上下文
	s.userService.BaseService.Ctx = ctx

	// 调用服务层
	user, err := s.userService.Create(createReq)
	if err != nil {
		return nil, convertError(err)
	}

	// 转换响应
	return &protobuf.CreateUserResponse{
		User: convertUserToProto(user),
	}, nil
}

// GetUser 获取用户
func (s *UserServer) GetUser(ctx context.Context, req *protobuf.GetUserRequest) (*protobuf.GetUserResponse, error) {
	// 设置上下文
	s.userService.BaseService.Ctx = ctx

	// 调用服务层
	user, err := s.userService.GetByID(uint(req.Id))
	if err != nil {
		return nil, convertError(err)
	}

	// 转换响应
	return &protobuf.GetUserResponse{
		User: convertUserToProto(user),
	}, nil
}

// UpdateUser 更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *protobuf.UpdateUserRequest) (*protobuf.UpdateUserResponse, error) {
	// 转换请求
	updateReq := &model.UpdateUserRequest{
		Name:   req.Name,
		Phone:  req.Phone,
		Avatar: req.Avatar,
	}

	// 设置上下文
	s.userService.BaseService.Ctx = ctx

	// 调用服务层
	user, err := s.userService.Update(uint(req.Id), updateReq)
	if err != nil {
		return nil, convertError(err)
	}

	// 转换响应
	return &protobuf.UpdateUserResponse{
		User: convertUserToProto(user),
	}, nil
}

// DeleteUser 删除用户
func (s *UserServer) DeleteUser(ctx context.Context, req *protobuf.DeleteUserRequest) (*protobuf.DeleteUserResponse, error) {
	// 设置上下文
	s.userService.BaseService.Ctx = ctx

	// 调用服务层
	err := s.userService.Delete(uint(req.Id))
	if err != nil {
		return nil, convertError(err)
	}

	// 返回响应
	return &protobuf.DeleteUserResponse{
		Success: true,
	}, nil
}

// ListUsers 获取用户列表
func (s *UserServer) ListUsers(ctx context.Context, req *protobuf.ListUsersRequest) (*protobuf.ListUsersResponse, error) {
	// 转换请求
	listReq := &model.UserListRequest{
		PageRequest: model.PageRequest{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		Username: req.Username,
		Email:    req.Email,
	}

	// 处理可选字段
	if req.Status != nil {
		status := int(*req.Status)
		listReq.Status = &status
	}

	// 设置默认值
	if listReq.Page == 0 {
		listReq.Page = 1
	}
	if listReq.Size == 0 {
		listReq.Size = 10
	}

	// 设置上下文
	s.userService.BaseService.Ctx = ctx

	// 调用服务层
	users, total, err := s.userService.List(listReq)
	if err != nil {
		return nil, convertError(err)
	}

	// 转换响应
	protoUsers := make([]*protobuf.User, len(users))
	for i, user := range users {
		protoUsers[i] = convertUserToProto(user)
	}

	return &protobuf.ListUsersResponse{
		Users: protoUsers,
		Total: total,
		Page:  int32(listReq.Page),
		Size:  int32(listReq.Size),
	}, nil
}

// Login 用户登录
func (s *UserServer) Login(ctx context.Context, req *protobuf.LoginRequest) (*protobuf.LoginResponse, error) {
	// 转换请求
	loginReq := &model.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	// 设置上下文
	s.userService.BaseService.Ctx = ctx

	// 调用服务层
	user, err := s.userService.Login(loginReq)
	if err != nil {
		return nil, convertError(err)
	}

	// 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	// 转换响应
	return &protobuf.LoginResponse{
		Token: token,
		User:  convertUserToProto(user),
	}, nil
}

// ChangePassword 修改密码
func (s *UserServer) ChangePassword(ctx context.Context, req *protobuf.ChangePasswordRequest) (*protobuf.ChangePasswordResponse, error) {
	// 设置上下文
	s.userService.BaseService.Ctx = ctx

	// 调用服务层
	err := s.userService.ChangePassword(uint(req.UserId), req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, convertError(err)
	}

	// 返回响应
	return &protobuf.ChangePasswordResponse{
		Success: true,
	}, nil
}

// convertUserToProto 转换 User 模型为 protobuf 格式
func convertUserToProto(user *model.User) *protobuf.User {
	protoUser := &protobuf.User{
		Id:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Name:       user.Name,
		Phone:      user.Phone,
		Avatar:     user.Avatar,
		Status:     int32(user.Status),
		LoginCount: int32(user.LoginCount),
		CreatedAt:  user.CreatedAt,  // 直接使用 int64 毫秒时间戳
		UpdatedAt:  user.UpdatedAt,  // 直接使用 int64 毫秒时间戳
	}

	// 处理可选字段
	if user.LastLoginAt != nil {
		protoUser.LastLoginAt = user.LastLoginAt.UnixMilli()  // 转换为毫秒时间戳
	}

	return protoUser
}

// convertError 转换错误为 gRPC 错误
func convertError(err error) error {
	if err == nil {
		return nil
	}

	// 检查是否是自定义错误
	if appErr, ok := err.(*errorsx.AppError); ok {
		switch appErr.Code {
		case errorsx.CodeUserNotFound:
			return status.Error(codes.NotFound, appErr.Message)
		case errorsx.CodeUserAlreadyExists:
			return status.Error(codes.AlreadyExists, appErr.Message)
		case errorsx.CodeInvalidCredentials:
			return status.Error(codes.Unauthenticated, appErr.Message)
		case errorsx.CodeUserDisabled:
			return status.Error(codes.PermissionDenied, appErr.Message)
		case errorsx.CodeValidationFailed:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case errorsx.CodeBadRequest:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case errorsx.CodeDatabaseError:
			return status.Error(codes.Internal, appErr.Message)
		case errorsx.CodeInternalServerError:
			return status.Error(codes.Internal, appErr.Message)
		default:
			return status.Error(codes.Internal, appErr.Message)
		}
	}

	// 默认返回内部错误
	return status.Error(codes.Internal, err.Error())
}
