#!/bin/bash

# Gin Starter Project Initialization Script
# 用于快速初始化新项目的脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认值
TEMPLATE_MODULE="github.com/iswangwenbin/gin-starter"
DEFAULT_PROJECT_NAME="my-gin-app"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 显示帮助信息
show_help() {
    cat << EOF
🚀 Gin Starter Project Initialization Tool

Usage: $0 [OPTIONS] <module-name> [project-directory]

Arguments:
  module-name         New module name (e.g., github.com/myorg/myapp)
  project-directory   Target directory (default: extracted from module name)

Options:
  -h, --help         Show this help message
  -f, --force        Force overwrite existing directory
  -i, --interactive  Interactive mode
  -v, --verbose      Verbose output

Examples:
  $0 github.com/myorg/awesome-api
  $0 -i github.com/company/project-name ./my-project
  $0 --force gitlab.com/team/service-api

EOF
}

# 解析命令行参数
parse_args() {
    FORCE=false
    INTERACTIVE=false
    VERBOSE=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -f|--force)
                FORCE=true
                shift
                ;;
            -i|--interactive)
                INTERACTIVE=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -*)
                print_error "Unknown option: $1"
                exit 1
                ;;
            *)
                if [ -z "$MODULE_NAME" ]; then
                    MODULE_NAME="$1"
                elif [ -z "$TARGET_DIR" ]; then
                    TARGET_DIR="$1"
                else
                    print_error "Too many arguments"
                    exit 1
                fi
                shift
                ;;
        esac
    done
}

# 验证模块名
validate_module_name() {
    if [ -z "$MODULE_NAME" ]; then
        print_error "Module name is required"
        show_help
        exit 1
    fi
    
    if [[ ! "$MODULE_NAME" =~ ^[a-zA-Z0-9._/-]+$ ]]; then
        print_error "Invalid module name format: $MODULE_NAME"
        exit 1
    fi
}

# 提取项目名
extract_project_name() {
    PROJECT_NAME=$(basename "$MODULE_NAME")
    if [ -z "$TARGET_DIR" ]; then
        TARGET_DIR="$PROJECT_NAME"
    fi
}

# 交互模式
interactive_mode() {
    echo -e "${BLUE}🎯 Interactive Project Setup${NC}"
    echo "==============================="
    
    read -p "Module name [$MODULE_NAME]: " input
    if [ -n "$input" ]; then
        MODULE_NAME="$input"
        PROJECT_NAME=$(basename "$MODULE_NAME")
    fi
    
    read -p "Project name [$PROJECT_NAME]: " input
    if [ -n "$input" ]; then
        PROJECT_NAME="$input"
    fi
    
    read -p "Target directory [$TARGET_DIR]: " input
    if [ -n "$input" ]; then
        TARGET_DIR="$input"
    fi
}

# 检查目标目录
check_target_directory() {
    if [ -d "$TARGET_DIR" ]; then
        if [ "$FORCE" = false ]; then
            print_error "Directory '$TARGET_DIR' already exists"
            print_info "Use --force to overwrite"
            exit 1
        else
            print_warning "Directory '$TARGET_DIR' exists, will overwrite..."
        fi
    fi
}

# 获取脚本所在目录（模板目录）
get_template_dir() {
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    TEMPLATE_DIR="$(dirname "$SCRIPT_DIR")"
}

# 复制文件并替换内容
copy_and_replace() {
    local src_dir="$1"
    local dst_dir="$2"
    
    # 需要排除的文件和目录
    local exclude_patterns=(
        ".git" ".idea" ".vscode" "*.log" "logs" "tmp" "bin" 
        "coverage.out" "coverage.html" ".DS_Store" "Thumbs.db"
        "scripts/init-project.sh"
    )
    
    # 创建目标目录
    mkdir -p "$dst_dir"
    
    # 使用 rsync 复制文件，排除指定模式
    local rsync_exclude=""
    for pattern in "${exclude_patterns[@]}"; do
        rsync_exclude="$rsync_exclude --exclude=$pattern"
    done
    
    if command -v rsync >/dev/null 2>&1; then
        eval "rsync -av $rsync_exclude '$src_dir/' '$dst_dir/'"
    else
        # 备用方案：使用 cp
        print_warning "rsync not found, using cp (less efficient)"
        cp -r "$src_dir"/* "$dst_dir/" 2>/dev/null || true
        
        # 手动删除排除的文件
        for pattern in "${exclude_patterns[@]}"; do
            find "$dst_dir" -name "$pattern" -type d -exec rm -rf {} + 2>/dev/null || true
            find "$dst_dir" -name "$pattern" -type f -delete 2>/dev/null || true
        done
    fi
}

# 替换文件内容
replace_module_name() {
    local target_dir="$1"
    
    print_info "Replacing module name in files..."
    
    # 需要处理的文件类型
    local file_patterns=(
        "*.go" "*.md" "*.txt" "*.yml" "*.yaml" "*.json" "*.toml"
        "*.sql" "*.sh" "*.dockerfile" "Makefile" "README*"
        ".gitignore" ".gitattributes" ".editorconfig"
    )
    
    for pattern in "${file_patterns[@]}"; do
        find "$target_dir" -name "$pattern" -type f -exec sed -i '' "s|$TEMPLATE_MODULE|$MODULE_NAME|g" {} + 2>/dev/null || true
        
        # 替换项目名称
        find "$target_dir" -name "$pattern" -type f -exec sed -i '' "s/gin-starter/$PROJECT_NAME/g" {} + 2>/dev/null || true
    done
    
    # 特殊处理 go.mod 文件
    if [ -f "$target_dir/go.mod" ]; then
        sed -i '' "1s/.*/module $MODULE_NAME/" "$target_dir/go.mod"
    fi
}

# 清理和初始化
cleanup_and_init() {
    local target_dir="$1"
    
    # 删除初始化脚本自身
    rm -f "$target_dir/scripts/init-project.sh"
    rm -f "$target_dir/cmd/init.go"
    
    # 重新生成 go.mod
    cd "$target_dir"
    if command -v go >/dev/null 2>&1; then
        print_info "Cleaning up Go modules..."
        go mod tidy
    else
        print_warning "Go not found, please run 'go mod tidy' manually"
    fi
    
    cd - >/dev/null
}

# 显示完成信息
show_completion_info() {
    echo
    print_success "Project initialized successfully!"
    echo
    echo -e "${BLUE}📋 Next steps:${NC}"
    echo "   cd $TARGET_DIR"
    
    if command -v make >/dev/null 2>&1; then
        echo "   make deps    # Install dependencies"
        echo "   make run     # Run the application"
    else
        echo "   go mod tidy  # Install dependencies"
        echo "   go run . serve  # Run the application"
    fi
    
    echo
    echo -e "${BLUE}📚 Useful commands:${NC}"
    echo "   make help    # Show all available commands"
    echo "   make test    # Run tests"
    echo "   make build   # Build the application"
    echo
}

# 主函数
main() {
    echo -e "${GREEN}🚀 Gin Starter Project Initialization Tool${NC}"
    echo "============================================="
    echo
    
    # 解析参数
    parse_args "$@"
    
    # 验证输入
    validate_module_name
    extract_project_name
    
    # 交互模式
    if [ "$INTERACTIVE" = true ]; then
        interactive_mode
    fi
    
    # 检查目标目录
    check_target_directory
    
    # 获取模板目录
    get_template_dir
    
    # 显示配置信息
    print_info "Module name: $MODULE_NAME"
    print_info "Project name: $PROJECT_NAME"
    print_info "Target directory: $TARGET_DIR"
    print_info "Template directory: $TEMPLATE_DIR"
    echo
    
    # 复制文件
    print_info "Copying template files..."
    copy_and_replace "$TEMPLATE_DIR" "$TARGET_DIR"
    
    # 替换模块名
    replace_module_name "$TARGET_DIR"
    
    # 清理和初始化
    cleanup_and_init "$TARGET_DIR"
    
    # 显示完成信息
    show_completion_info
}

# 运行主函数
main "$@"