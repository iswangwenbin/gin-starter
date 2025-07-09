# Makefile for gin-starter project

# 变量定义
APP_NAME := gin-starter
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version)

# 构建标志
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X "main.GoVersion=$(GO_VERSION)"

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "可用的命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 构建相关
.PHONY: build
build: ## 构建应用程序
	@echo "Building $(APP_NAME)..."
	@go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) .

.PHONY: build-linux
build-linux: ## 构建 Linux 版本
	@echo "Building $(APP_NAME) for Linux..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-linux .

.PHONY: build-windows
build-windows: ## 构建 Windows 版本
	@echo "Building $(APP_NAME) for Windows..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-windows.exe .

.PHONY: build-all
build-all: build build-linux build-windows ## 构建所有平台版本

# 测试相关
.PHONY: test
test: ## 运行测试
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

.PHONY: test-coverage
test-coverage: test ## 运行测试并生成覆盖率报告
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-integration
test-integration: ## 运行集成测试
	@echo "Running integration tests..."
	@go test -v -tags=integration ./tests/...

.PHONY: benchmark
benchmark: ## 运行基准测试
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# 代码质量
.PHONY: lint
lint: ## 运行代码检查
	@echo "Running linter..."
	@golangci-lint run

.PHONY: fmt
fmt: ## 格式化代码
	@echo "Formatting code..."
	@go fmt ./...

.PHONY: vet
vet: ## 运行 go vet
	@echo "Running go vet..."
	@go vet ./...

.PHONY: tidy
tidy: ## 整理依赖
	@echo "Tidying dependencies..."
	@go mod tidy

.PHONY: verify
verify: fmt vet lint test ## 验证代码质量（格式化、检查、测试）

# 运行相关
.PHONY: run
run: ## 运行应用程序
	@echo "Running $(APP_NAME)..."
	@go run . serve

.PHONY: run-dev
run-dev: ## 以开发模式运行
	@echo "Running $(APP_NAME) in development mode..."
	@go run . serve --env development --debug

.PHONY: run-prod
run-prod: ## 以生产模式运行
	@echo "Running $(APP_NAME) in production mode..."
	@go run . serve --env production

# Docker 相关
.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):$(VERSION) .
	@docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 $(APP_NAME):latest

.PHONY: docker-push
docker-push: ## 推送 Docker 镜像
	@echo "Pushing Docker image..."
	@docker push $(APP_NAME):$(VERSION)
	@docker push $(APP_NAME):latest

# 开发工具
.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/google/wire/cmd/wire@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: generate
generate: ## 生成代码
	@echo "Generating code..."
	@go generate ./...

.PHONY: proto
proto: ## 生成 protobuf 代码
	@echo "Generating protobuf code..."
	@./scripts/generate-proto.sh

.PHONY: swagger
swagger: ## 生成 Swagger 文档
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/main.go -o ./docs

.PHONY: wire
wire: ## 生成依赖注入代码
	@echo "Generating wire code..."
	@wire ./...

# 数据库相关
.PHONY: db-migrate
db-migrate: ## 运行数据库迁移
	@echo "Running database migrations..."
	@go run . migrate

.PHONY: db-seed
db-seed: ## 填充测试数据
	@echo "Seeding database..."
	@go run . seed

# 清理
.PHONY: clean
clean: ## 清理构建文件
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean -cache

.PHONY: deps
deps: ## 下载依赖
	@echo "Downloading dependencies..."
	@go mod download

# 发布相关
.PHONY: release
release: verify build-all ## 准备发布（验证 + 构建所有平台）
	@echo "Release ready!"

# 监控文件变化并自动重新构建
.PHONY: watch
watch: ## 监控文件变化并自动重新构建
	@echo "Watching for changes..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

# 安全检查
.PHONY: security
security: ## 运行安全检查
	@echo "Running security checks..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	@gosec ./...

# 性能分析
.PHONY: profile
profile: ## 运行性能分析
	@echo "Running performance profiling..."
	@go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...
	@echo "Profile files generated: cpu.prof, mem.prof"

# 全面检查
.PHONY: check-all
check-all: deps verify security ## 全面检查（依赖、验证、安全）

# 本地开发环境初始化
.PHONY: init
init: deps install-tools ## 初始化开发环境
	@echo "Development environment initialized!"