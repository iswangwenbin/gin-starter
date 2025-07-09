# 单进程双协议架构说明

## 📋 架构概述

本项目采用**单进程双协议**架构，在一个进程中同时运行 HTTP 和 gRPC 服务，实现了以下优势：

- ✅ **资源共享高效**: 数据库连接池、Redis连接、内存缓存等资源共享
- ✅ **配置管理统一**: 一套配置文件管理所有服务
- ✅ **部署运维简单**: 只需部署一个二进制文件
- ✅ **代码复用度高**: Service、Repository 层完全共享
- ✅ **监控告警简单**: 只需监控一个进程

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────┐
│                 单进程应用                           │
├─────────────────────────────────────────────────────┤
│  HTTP Server (Port: 9000)   │  gRPC Server (Port: 9090) │
│  ┌─────────────────────────┐ │  ┌─────────────────────────┐ │
│  │   Gin Router            │ │  │   gRPC Services         │ │
│  │   - REST API            │ │  │   - UserService         │ │
│  │   - Middleware          │ │  │   - HealthService       │ │
│  │   - Authentication     │ │  │   - Interceptors        │ │
│  └─────────────────────────┘ │  └─────────────────────────┘ │
├─────────────────────────────────────────────────────┤
│                 共享业务层                           │
│  ┌─────────────────────────────────────────────────┐ │
│  │             Service Layer                       │ │
│  │  - UserService (业务逻辑)                       │ │
│  │  - 数据验证、权限检查                            │ │
│  │  - 事务处理                                      │ │
│  └─────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────┐ │
│  │            Repository Layer                     │ │
│  │  - UserRepository (数据访问)                    │ │
│  │  - 数据库操作封装                                │ │
│  │  - 缓存操作                                      │ │
│  └─────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────┤
│                 共享资源层                           │
│  ┌─────────────────────────────────────────────────┐ │
│  │  Database Pool  │  Redis Pool  │  Logger        │ │
│  │  Configuration  │  Monitoring  │  Caching       │ │
│  └─────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

## ⚙️ 配置管理

```yaml
# config/local.yaml
server:
  host: 0.0.0.0
  port: "9000"        # HTTP 端口
  
grpc:
  port: 9090          # gRPC 端口
  enabled: true       # 是否启用 gRPC

debug: true           # 调试模式
```

## 🚀 启动方式

### 默认启动（HTTP + gRPC）
```bash
# 启动双协议服务
go run . serve

# 或者使用编译后的二进制
./bin/gin-starter serve
```

### 选择性启动

```bash
# 仅启动 HTTP 服务
go run . serve --env http-only

# 仅启动 gRPC 服务
go run . serve --env grpc-only
```

## 📊 服务状态监控

```bash
# 检查 HTTP 服务
curl http://localhost:9000/health

# 检查 gRPC 服务
grpcurl -plaintext localhost:9090 common.HealthService/Check
```

## 🔧 开发模式

### 实时监控文件变化
```bash
make watch   # 自动重新编译和重启
```

### 并发测试
```bash
# 同时测试 HTTP 和 gRPC
go run examples/concurrent_test.go
```

## 📝 日志示例

启动时会看到类似以下日志：

```
2024-01-01T10:00:00Z INFO HTTP server starting address=:9000
2024-01-01T10:00:00Z INFO gRPC server starting address=:9090
2024-01-01T10:00:00Z INFO Servers started successfully http_address=:9000 grpc_address=:9090 grpc_enabled=true
```

## 🛠️ 开发指南

### 添加新的 API 端点

1. **HTTP**: 在 `internal/api/` 添加路由和处理函数
2. **gRPC**: 在 `internal/grpc/protobuf/` 添加 proto 定义，运行 `make proto` 生成代码

### 共享业务逻辑

两种协议的服务都通过相同的 Service 层处理业务逻辑：

```go
// HTTP 控制器
func (uc *UserController) Create(c *gin.Context) {
    // 解析 HTTP 请求
    var req model.CreateUserRequest
    c.ShouldBindJSON(&req)
    
    // 调用共享的业务逻辑
    user, err := uc.userService.Create(&req)
    // ...
}

// gRPC 服务
func (s *UserServer) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.CreateUserResponse, error) {
    // 转换 gRPC 请求
    createReq := &model.CreateUserRequest{
        Username: req.Username,
        Email:    req.Email,
        // ...
    }
    
    // 调用相同的业务逻辑
    user, err := s.userService.Create(createReq)
    // ...
}
```

## 🔍 故障排除

### 端口冲突
```bash
# 检查端口占用
lsof -i :9000
lsof -i :9090

# 修改配置文件中的端口
```

### 服务启动失败
```bash
# 检查日志
tail -f logs/server.log

# 查看详细错误
go run . serve --debug
```

## 📈 性能优化

### 资源池配置
```yaml
database:
  max_idle_conns: 10
  max_open_conns: 100
  
redis:
  pool_size: 10
```

### 并发处理
- HTTP: Gin 框架自动处理并发
- gRPC: 使用 goroutine 处理每个请求

## 🔒 安全考虑

- 共享认证机制
- 统一的错误处理
- 相同的安全中间件
- 统一的日志记录

## 📊 监控指标

系统提供统一的监控指标：

- 请求数量（HTTP + gRPC）
- 响应时间分布
- 错误率统计
- 资源使用情况

## 🤔 何时考虑拆分

只有在以下情况下才考虑拆分为两个进程：

1. **负载特征差异巨大**: HTTP 主要外部调用，gRPC 主要内部高频调用
2. **扩展需求不同**: 需要独立扩展不同协议的服务
3. **技术栈要求不同**: 需要不同的中间件或框架
4. **团队分工明确**: 不同团队负责不同协议

## 🎯 最佳实践

1. **统一错误处理**: 在 Service 层统一错误格式
2. **共享数据模型**: 使用相同的数据模型定义
3. **一致的日志格式**: 统一的日志记录和追踪
4. **配置中心化**: 所有配置集中管理
5. **监控统一化**: 使用相同的监控工具和指标

这种架构既保持了系统的简洁性，又提供了足够的灵活性来应对不同的客户端需求。