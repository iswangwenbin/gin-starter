# gRPC 架构说明

本项目采用清晰的分层架构来组织 gRPC 相关代码：

## 目录结构

```
internal/grpc/
├── protobuf/           # Protocol Buffers 定义和生成的代码
│   ├── *.proto        # 协议定义文件
│   ├── *.pb.go        # 生成的 Go 结构体
│   └── *_grpc.pb.go   # 生成的 gRPC 服务代码
├── server/            # gRPC 服务端实现
│   ├── server.go      # 服务器管理器
│   ├── user_server.go # 用户服务实现
│   └── health_server.go # 健康检查服务实现
└── client/            # gRPC 客户端代码
    ├── client.go      # 客户端管理器
    └── user_client.go # 用户服务客户端
```

## Protocol Buffers

### 文件说明

- `user.proto`: 用户服务的协议定义
- `common.proto`: 通用消息类型和健康检查服务

### 代码生成

```bash
# 生成 protobuf 代码
make proto

# 或者直接运行脚本
./scripts/generate-proto.sh
```

## 服务端实现

### 启动 gRPC 服务器

```go
package main

import (
    "context"
    "log"
    
    "github.com/iswangwenbin/gin-starter/internal/grpc/server"
    "github.com/iswangwenbin/gin-starter/pkg/configx"
)

func main() {
    // 加载配置
    config, err := configx.Load("config/local.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建服务器
    grpcServer := server.NewServer(config, logger, db, cache)
    
    // 启动服务器
    ctx := context.Background()
    if err := grpcServer.Start(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### 配置

在 `config/local.yaml` 中添加：

```yaml
grpc:
  port: 9090
  enabled: true

debug: true  # 启用 gRPC 反射
```

## 客户端使用

### 创建客户端

```go
package main

import (
    "context"
    "log"
    
    "github.com/iswangwenbin/gin-starter/internal/grpc/client"
    "github.com/iswangwenbin/gin-starter/internal/grpc/protobuf"
)

func main() {
    // 创建客户端管理器
    clientManager, err := client.NewClientManager("localhost:9090", logger)
    if err != nil {
        log.Fatal(err)
    }
    defer clientManager.Close()
    
    // 获取用户客户端
    userClient := clientManager.UserClient()
    
    // 调用服务
    req := &protobuf.CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    resp, err := userClient.CreateUser(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("User created: %v", resp.User)
}
```

## 服务列表

### UserService

用户管理服务，提供以下 RPC 方法：

- `CreateUser`: 创建用户
- `GetUser`: 获取用户信息
- `UpdateUser`: 更新用户信息
- `DeleteUser`: 删除用户
- `ListUsers`: 获取用户列表
- `Login`: 用户登录
- `ChangePassword`: 修改密码

### HealthService

健康检查服务：

- `Check`: 检查服务健康状态

## 错误处理

gRPC 服务端会自动将应用错误转换为对应的 gRPC 状态码：

- `CodeUserNotFound` → `codes.NotFound`
- `CodeUserAlreadyExists` → `codes.AlreadyExists`
- `CodeInvalidCredentials` → `codes.Unauthenticated`
- `CodeValidationFailed` → `codes.InvalidArgument`
- `CodeDatabaseError` → `codes.Internal`

## 测试

### 使用 grpcurl 测试

```bash
# 安装 grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 列出服务
grpcurl -plaintext localhost:9090 list

# 调用健康检查
grpcurl -plaintext localhost:9090 common.HealthService/Check

# 创建用户
grpcurl -plaintext -d '{
  "username": "testuser",
  "email": "test@example.com", 
  "password": "password123",
  "name": "Test User"
}' localhost:9090 user.UserService/CreateUser
```

### 运行示例客户端

```bash
go run examples/grpc_client_example.go
```

## 开发工具

### 安装 Protocol Buffers 编译器

```bash
# macOS
brew install protobuf

# 安装 Go 插件
make install-tools
```

### 生成代码

```bash
# 生成 protobuf 代码
make proto
```

### 构建和运行

```bash
# 构建项目
make build

# 运行服务器（包含 gRPC）
./bin/gin-starter serve
```

## 与 HTTP API 的关系

gRPC 服务与现有的 HTTP API 并行运行：

- HTTP API 端口：9000（可配置）
- gRPC 端口：9090（可配置）

两者共享相同的业务逻辑层（Service）和数据访问层（Repository），确保数据一致性。

## 最佳实践

1. **协议版本管理**: 使用向后兼容的方式修改 .proto 文件
2. **错误处理**: 统一的错误转换机制
3. **日志记录**: 所有 gRPC 调用都会记录详细日志
4. **健康检查**: 实现标准的健康检查接口
5. **连接池**: 客户端使用连接池管理连接
6. **超时控制**: 为所有 RPC 调用设置合理的超时时间