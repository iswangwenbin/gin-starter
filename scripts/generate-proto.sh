#!/bin/bash

# gRPC 代码生成脚本

set -e

# 检查 protoc 是否安装
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc is not installed. Please install Protocol Buffers compiler."
    echo "   brew install protobuf"
    echo "   or download from https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# 检查 protoc-gen-go 是否安装
if ! command -v protoc-gen-go &> /dev/null; then
    echo "❌ protoc-gen-go is not installed. Installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# 检查 protoc-gen-go-grpc 是否安装
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "❌ protoc-gen-go-grpc is not installed. Installing..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROTO_DIR="$PROJECT_ROOT/internal/grpc/protobuf"

echo "🚀 Generating gRPC code..."
echo "📁 Project root: $PROJECT_ROOT"
echo "📁 Proto directory: $PROTO_DIR"

# 切换到项目根目录
cd "$PROJECT_ROOT"

# 创建输出目录
mkdir -p "$PROTO_DIR"

# 生成 protobuf 文件
for proto_file in "$PROTO_DIR"/*.proto; do
    if [ -f "$proto_file" ]; then
        echo "🔧 Processing $(basename "$proto_file")..."
        
        protoc \
            --proto_path="$PROTO_DIR" \
            --go_out="$PROTO_DIR" \
            --go_opt=paths=source_relative \
            --go-grpc_out="$PROTO_DIR" \
            --go-grpc_opt=paths=source_relative \
            "$proto_file"
    fi
done

echo "✅ gRPC code generation completed!"
echo "📋 Generated files:"
find "$PROTO_DIR" -name "*.pb.go" -o -name "*_grpc.pb.go" | sort