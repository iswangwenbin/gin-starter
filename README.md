# Gin Starter

一个基于 Gin 框架的 Go 微服务起始项目，提供完整的企业级功能和最佳实践。

## 🚀 特性

### 核心功能
- **Web 框架**: 基于 [Gin](https://gin-gonic.com/) 的高性能 HTTP 服务
- **gRPC 支持**: 内置 gRPC 服务器，支持 HTTP 和 gRPC 双协议
- **数据库**: MySQL + GORM ORM，支持数据库迁移
- **缓存**: Redis 集成，支持分布式缓存
- **时序数据**: ClickHouse 支持，用于日志和事件数据存储
- **消息队列**: Redis Stream 异步事件处理

### 架构特性
- **分层架构**: API → Service → Repository → Model 清晰分层
- **依赖注入**: 基于接口的依赖注入设计
- **配置管理**: 多环境配置支持（development/production/local）
- **日志系统**: 基于 Zap 的结构化日志
- **错误处理**: 统一的错误码和错误处理机制

### 开发特性
- **JWT 认证**: 完整的用户认证和授权系统
- **中间件**: CORS、限流、安全头、请求追踪等
- **Docker 支持**: 完整的 Docker Compose 开发环境
- **代码生成**: Protocol Buffers 代码生成
- **命令行工具**: 基于 Cobra 的 CLI 工具

## 📁 项目结构

```
gin-starter/
├── cmd/                    # 命令行工具
│   ├── root.go            # 根命令和全局配置
│   ├── serve.go           # HTTP/gRPC 服务器
│   └── worker.go          # 异步事件处理器
├── config/                # 配置文件
│   ├── development.yaml   # 开发环境配置
│   ├── production.yaml    # 生产环境配置
│   └── local.yaml         # 本地开发配置
├── internal/              # 内部包（不对外暴露）
│   ├── api/               # HTTP API 处理器
│   ├── core/              # 核心服务器逻辑
│   ├── grpc/              # gRPC 服务实现
│   ├── middleware/        # HTTP 中间件
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   ├── service/           # 业务逻辑层
│   └── worker/            # 后台工作器
├── pkg/                   # 公共包（可对外暴露）
│   ├── configx/           # 配置工具
│   ├── databasex/         # 数据库工具
│   ├── clickhousex/       # ClickHouse 工具
│   ├── redisx/            # Redis 工具
│   └── errorsx/           # 错误处理工具
└── docker/                # Docker 相关文件
    └── docker-compose.yml # 开发环境服务编排
```

## 🔧 快速开始

### 环境要求
- Go 1.21+
- Docker & Docker Compose
- Protocol Buffers 编译器（可选，用于 gRPC 开发）

### 1. 克隆项目
```bash
git clone <repository-url>
cd gin-starter
```

### 2. 启动开发环境
```bash
# 启动数据库和中间件服务
docker-compose -f docker/docker-compose.yml up -d

# 安装依赖
go mod tidy
```

### 3. 运行服务
```bash
# 启动 HTTP/gRPC 服务器
go run main.go serve --env local --debug

# 或者启动异步事件处理器
go run main.go worker --env local --debug
```

### 4. 验证服务
```bash
# 健康检查
curl http://localhost:8001/health

# API 测试
curl http://localhost:8001/api/v1/ping
```

## 🛠️ 开发指南

### 命令行工具

项目提供了完整的 CLI 工具：

```bash
# 查看所有可用命令
go run main.go --help

# 启动 HTTP/gRPC 服务器
go run main.go serve [flags]

# 启动事件处理器
go run main.go worker [flags]

# 初始化新项目
go run main.go init [flags]
```

#### 全局标志
- `--env string`: 运行环境 (development/production/local)
- `--debug`: 启用调试模式
- `--config string`: 指定配置文件路径

### 配置管理

项目支持多环境配置，配置文件位于 `config/` 目录：

- `local.yaml`: 本地开发环境
- `development.yaml`: 开发环境
- `production.yaml`: 生产环境

### 数据库迁移

使用 GORM 的自动迁移功能：

```go
// 在 internal/core/server.go 中
db.AutoMigrate(&model.User{}, &model.InstallEvent{})
```

### gRPC 开发

1. 编辑 `.proto` 文件：`internal/grpc/protobuf/`
2. 生成代码：`./scripts/generate-proto.sh`
3. 实现服务：`internal/grpc/server/`

### 添加新的 API

1. 定义模型：`internal/model/`
2. 实现 Repository：`internal/repository/`
3. 实现 Service：`internal/service/`
4. 实现 API Handler：`internal/api/`
5. 添加路由：`internal/core/routes.go`

## 🏗️ 架构设计

### 分层架构

```
┌─────────────────┐
│   HTTP/gRPC     │  ← API 层（路由、中间件、参数验证）
├─────────────────┤
│    Service      │  ← 业务逻辑层（业务规则、事务处理）
├─────────────────┤
│   Repository    │  ← 数据访问层（数据库、缓存操作）
├─────────────────┤
│     Model       │  ← 数据模型层（实体定义）
└─────────────────┘
```

### 异步处理架构

```
┌──────────┐    ┌─────────────────┐    ┌─────────────┐
│ HTTP API │───▶│  Redis Stream   │───▶│   Worker    │
└──────────┘    └─────────────────┘    └─────────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │   ClickHouse    │
                                    └─────────────────┘
```

### 进程模型

- **serve**: HTTP/gRPC 服务器，处理 API 请求
- **worker**: 异步事件处理器，消费消息队列

## 📦 依赖管理

### 核心依赖

```bash
# Web 框架
go get -u github.com/gin-gonic/gin

# 数据库 ORM
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

# Redis 客户端
go get github.com/redis/go-redis/v9

# ClickHouse 客户端
go get -u github.com/ClickHouse/clickhouse-go/v2

# gRPC
go get -u google.golang.org/grpc
go get -u google.golang.org/protobuf

# 日志
go get -u go.uber.org/zap
go get -u github.com/gin-contrib/zap

# JWT
go get -u github.com/golang-jwt/jwt/v5

# 命令行工具
go get -u github.com/spf13/cobra

# 工具库
go get -u github.com/spf13/cast
go get -u github.com/pkg/errors
```

## 🚀 部署

### Docker 部署

```bash
# 构建镜像
docker build -t gin-starter .

# 运行容器
docker run -p 8001:8001 -p 8002:8002 gin-starter
```

### 环境变量

生产环境可以通过环境变量覆盖配置：

```bash
export JWT_SECRET="your-production-secret"
export REDIS_PASSWORD="redis-password"
export CLICKHOUSE_HOST="clickhouse-host"
export CLICKHOUSE_PASSWORD="clickhouse-password"
```

## 📊 监控和日志

- **结构化日志**: 使用 Zap 记录结构化日志
- **请求追踪**: 每个请求都有唯一的 Request ID
- **性能监控**: 内置性能指标收集
- **健康检查**: `/health` 端点提供服务健康状态

## 🤝 贡献

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交改动 (`git commit -m 'Add some amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关链接

- [Gin 文档](https://gin-gonic.com/)
- [GORM 文档](https://gorm.io/)
- [Redis 文档](https://redis.io/)
- [ClickHouse 文档](https://clickhouse.com/)
- [gRPC Go 教程](https://grpc.io/docs/languages/go/)