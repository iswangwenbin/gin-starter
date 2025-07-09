# Gin Starter - 企业级 Web 框架模版

这是一个基于 Gin 框架构建的完整企业级 Web 应用模版，集成了现代 Web 开发所需的各种功能和最佳实践。

## 🚀 主要特性

### 核心功能

- **RESTful API 架构** - 清晰的 API 设计和路由管理
- **JWT 认证系统** - 基于 JSON Web Token 的用户认证
- **用户管理系统** - 完整的用户注册、登录、个人信息管理
- **配置管理** - 多环境配置支持（开发、生产、本地）
- **数据库集成** - GORM ORM 与 MySQL 数据库
- **Redis 缓存** - Redis 缓存支持和连接池管理
- **健康检查** - 应用和依赖服务的健康状态监控

### 中间件系统

- **CORS 跨域处理** - 可配置的跨域资源共享
- **限流控制** - 基于 Redis 的 IP 和用户限流
- **请求 ID 追踪** - 每个请求的唯一标识符
- **安全头设置** - 常见的 Web 安全头配置
- **错误恢复** - 优雅的 panic 恢复和错误处理
- **日志记录** - 基于 Zap 的高性能日志系统

### 数据验证

- **自定义验证器** - 手机号、用户名、密码强度验证
- **统一错误响应** - 标准化的 API 错误响应格式
- **模型绑定** - 自动的 JSON/Query 参数绑定和验证

## 📁 项目结构

┌─────────────┐
│ API │ ← 控制器层（处理 HTTP 请求）
├─────────────┤
│ Service │ ← 业务逻辑层（业务规则和流程）
├─────────────┤
│ Repository │ ← 数据访问层（数据库操作抽象）
├─────────────┤
│ Model │ ← 数据模型层（数据结构定义）
└─────────────┘

```
gin-starter/
├── cmd/                    # 命令行工具
│   ├── root.go            # 根命令
│   └── serve.go           # 服务启动命令
├── config/                # 配置文件
│   ├── development.yaml   # 开发环境配置
│   ├── production.yaml    # 生产环境配置
│   └── local.yaml         # 本地配置
├── internel/              # 内部业务逻辑
│   ├── api/               # API 控制器和响应处理
│   │   ├── base.go        # 基础控制器
│   │   ├── response.go    # 统一响应格式
│   │   ├── error.go       # 错误处理
│   │   ├── validator.go   # 数据验证
│   │   ├── health.go      # 健康检查
│   │   └── user.go        # 用户管理
│   ├── middleware/        # 中间件
│   │   ├── cors.go        # 跨域处理
│   │   ├── jwt.go         # JWT 认证
│   │   ├── rate_limit.go  # 限流控制
│   │   ├── request_id.go  # 请求 ID
│   │   ├── security.go    # 安全头
│   │   └── error_handler.go # 错误处理
│   ├── model/             # 数据模型
│   │   ├── base.go        # 基础模型
│   │   └── user.go        # 用户模型
│   ├── service/           # 业务逻辑层
│   │   ├── base.go        # 基础服务
│   │   └── user.go        # 用户服务
│   ├── options.go         # 服务器选项
│   ├── routes.go          # 路由配置
│   └── server.go          # 服务器核心
├── pkg/                   # 公共包
│   ├── config/            # 配置管理
│   ├── databasex/         # 数据库连接
│   ├── logx/              # 日志工具
│   └── redisx/            # Redis 连接
├── docker/                # Docker 相关文件
└── main.go                # 应用入口
```

## 🛠️ 技术栈

- **Web 框架**: Gin
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **认证**: JWT
- **日志**: Zap
- **配置**: Viper
- **CLI**: Cobra
- **容器**: Docker

## 🔧 安装和使用

### 环境要求

- Go 1.21+
- MySQL 5.7+
- Redis 6.0+

### 1. 克隆项目

```bash
git clone <repository-url>
cd gin-starter
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库和 Redis

根据你的环境修改 `config/` 目录下的配置文件。

### 4. 启动服务

```bash
# 开发环境
go run main.go serve --env development

# 生产环境
go run main.go serve --env production

# 或构建后运行
go build -o gin-starter
./gin-starter serve --env development
```

## 📚 API 文档

### 认证相关

```
POST /api/v1/auth/register  # 用户注册
POST /api/v1/auth/login     # 用户登录
```

### 用户管理（需要认证）

```
GET    /api/v1/profile           # 获取个人信息
PUT    /api/v1/profile           # 更新个人信息
POST   /api/v1/change-password   # 修改密码
GET    /api/v1/users             # 用户列表（管理员）
GET    /api/v1/users/:id         # 获取用户信息
PUT    /api/v1/users/:id         # 更新用户信息
DELETE /api/v1/users/:id         # 删除用户
```

### 健康检查

```
GET /ping                # 简单 ping
GET /health              # 详细健康检查
GET /api/v1/ping         # API ping
GET /api/v1/health       # API 健康检查
```

## 🎯 配置说明

### 环境变量

生产环境建议使用环境变量覆盖配置：

```bash
export DB_PASSWORD="your-db-password"
export REDIS_PASSWORD="your-redis-password"
export JWT_SECRET="your-jwt-secret"
export FRONTEND_URL="https://your-frontend.com"
```

### 配置文件示例

```yaml
server:
  host: 0.0.0.0
  port: "8080"
  mode: release

database:
  host: localhost
  port: 3306
  user: root
  password: "${DB_PASSWORD}"
  name: gin_starter

jwt:
  secret: "${JWT_SECRET}"
  expires: 24h

cors:
  allowed_origins:
    - "${FRONTEND_URL}"
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
  allow_credentials: true
```

## 🔒 安全特性

- **JWT Token 认证** - 安全的用户认证机制
- **密码加密** - 使用 bcrypt 哈希算法
- **CORS 保护** - 可配置的跨域访问控制
- **限流保护** - 防止恶意请求和 DDoS 攻击
- **安全头** - 设置常见的 Web 安全响应头
- **输入验证** - 严格的数据验证和清理

## 🚀 部署

### Docker 部署

项目已包含 Docker 配置，可以使用 Docker Compose 进行部署：

```bash
docker-compose up -d
```

### 生产环境建议

- 使用反向代理（Nginx）
- 启用 HTTPS
- 配置日志轮转
- 设置监控告警
- 定期备份数据

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
