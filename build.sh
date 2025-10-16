#!/bin/bash

# 小红书 MCP 项目编译脚本
# 支持选择性编译: home/login/mcp/all

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示使用说明
show_usage() {
    echo "使用方法: $0 [target]"
    echo ""
    echo "编译目标:"
    echo "  home    - 编译首页交互工具"
    echo "  login   - 编译登录工具"
    echo "  mcp     - 编译 MCP 服务器"
    echo "  all     - 编译所有程序 (默认)"
    echo ""
    echo "示例:"
    echo "  $0 home   # 仅编译 home 工具"
    echo "  $0 all    # 编译所有程序"
    echo "  $0        # 编译所有程序 (默认)"
}

# 获取项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")" && pwd)
BUILD_DIR="$PROJECT_ROOT/build"

# 解析命令行参数
TARGET="${1:-all}"

case "$TARGET" in
    home|login|mcp|all)
        ;;
    -h|--help|help)
        show_usage
        exit 0
        ;;
    *)
        print_error "未知的编译目标: $TARGET"
        show_usage
        exit 1
        ;;
esac

print_info "开始编译小红书 MCP 项目..."
print_info "编译目标: $TARGET"
print_info "项目根目录: $PROJECT_ROOT"
print_info "编译输出目录: $BUILD_DIR"

# 创建 build 目录
if [ -d "$BUILD_DIR" ]; then
    print_warning "build 目录已存在，清理旧文件..."
    rm -rf "$BUILD_DIR"
fi

mkdir -p "$BUILD_DIR"
print_success "创建 build 目录: $BUILD_DIR"

# 进入项目根目录
cd "$PROJECT_ROOT"

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    print_error "Go 未安装或不在 PATH 中"
    exit 1
fi

GO_VERSION=$(go version)
print_info "Go 版本: $GO_VERSION"

# 编译函数
compile_mcp() {
    print_info "编译 MCP 服务器..."
    go build -o "$BUILD_DIR/mcp" .
    if [ $? -eq 0 ]; then
        print_success "MCP 服务器编译完成: $BUILD_DIR/mcp"
    else
        print_error "MCP 服务器编译失败"
        exit 1
    fi
}

compile_login() {
    print_info "编译登录工具..."
    go build -o "$BUILD_DIR/login" ./cmd/login
    if [ $? -eq 0 ]; then
        print_success "登录工具编译完成: $BUILD_DIR/login"
    else
        print_error "登录工具编译失败"
        exit 1
    fi
}

compile_home() {
    print_info "编译首页交互工具..."
    go build -o "$BUILD_DIR/home" ./cmd/home
    if [ $? -eq 0 ]; then
        print_success "首页交互工具编译完成: $BUILD_DIR/home"
    else
        print_error "首页交互工具编译失败"
        exit 1
    fi
}

# 根据目标进行编译
case "$TARGET" in
    mcp)
        compile_mcp
        ;;
    login)
        compile_login
        ;;
    home)
        compile_home
        ;;
    all)
        compile_mcp
        compile_login
        compile_home
        ;;
esac

# 显示编译结果
print_info "编译完成！文件列表："
ls -la "$BUILD_DIR"

# 显示文件大小
print_info "编译产物大小："
du -h "$BUILD_DIR"/*

print_success "编译任务完成！"
print_info "编译产物位置: $BUILD_DIR"

# 根据编译目标给出使用提示
case "$TARGET" in
    mcp)
        print_info "运行 './build/mcp' 启动 MCP 服务器"
        ;;
    login)
        print_info "运行 './build/login' 进行登录"
        ;;
    home)
        print_info "运行 './build/home' 进行首页交互"
        ;;
    all)
        print_info "运行 './build/login' 开始使用"
        ;;
esac