#!/bin/bash

# 小红书 MCP 项目编译脚本
# 编译所有可执行文件并放置到 build 目录

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

# 获取项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")" && pwd)
BUILD_DIR="$PROJECT_ROOT/build"

print_info "开始编译小红书 MCP 项目..."
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

# 编译主程序 (MCP 服务器)
print_info "编译主程序 (MCP 服务器)..."
go build -o "$BUILD_DIR/xiaohongshu-mcp" .
if [ $? -eq 0 ]; then
    print_success "主程序编译完成: $BUILD_DIR/xiaohongshu-mcp"
else
    print_error "主程序编译失败"
    exit 1
fi

# 编译登录工具
print_info "编译登录工具..."
go build -o "$BUILD_DIR/login" ./cmd/login
if [ $? -eq 0 ]; then
    print_success "登录工具编译完成: $BUILD_DIR/login"
else
    print_error "登录工具编译失败"
    exit 1
fi

# 编译首页交互工具
print_info "编译首页交互工具..."
go build -o "$BUILD_DIR/home" ./cmd/home
if [ $? -eq 0 ]; then
    print_success "首页交互工具编译完成: $BUILD_DIR/home"
else
    print_error "首页交互工具编译失败"
    exit 1
fi

# 复制配置文件和文档
print_info "复制配置文件和文档..."

# 复制 README 文件
cp README.md "$BUILD_DIR/" 2>/dev/null || print_warning "README.md 不存在"
cp README_EN.md "$BUILD_DIR/" 2>/dev/null || print_warning "README_EN.md 不存在"

# 复制 cmd 目录的 README
cp cmd/README.md "$BUILD_DIR/cmd_README.md" 2>/dev/null || print_warning "cmd/README.md 不存在"

# 复制示例配置文件
mkdir -p "$BUILD_DIR/examples"
if [ -d "examples" ]; then
    cp -r examples/* "$BUILD_DIR/examples/" 2>/dev/null || print_warning "复制示例文件失败"
fi

# 创建使用说明
cat > "$BUILD_DIR/USAGE.md" << 'EOF'
# 小红书 MCP 使用说明

## 编译产物说明

- `xiaohongshu-mcp`: 主程序，MCP 服务器
- `login`: 登录工具，用于获取和保存登录状态
- `home`: 首页交互工具，用于浏览器交互

## 使用步骤

1. 首先运行登录工具获取登录状态：
   ```bash
   ./login
   ```

2. 登录成功后，可以运行首页交互工具：
   ```bash
   ./home
   ```

3. 或者启动 MCP 服务器：
   ```bash
   ./xiaohongshu-mcp
   ```

## 注意事项

- 确保系统已安装 Chrome 或 Chromium 浏览器
- 登录状态会保存在 cookies.json 文件中
- 首次使用需要先运行 login 工具

更多详细信息请参考 README.md 文件。
EOF

print_success "使用说明创建完成: $BUILD_DIR/USAGE.md"

# 显示编译结果
print_info "编译完成！文件列表："
ls -la "$BUILD_DIR"

# 显示文件大小
print_info "编译产物大小："
du -h "$BUILD_DIR"/*

print_success "所有编译任务完成！"
print_info "编译产物位置: $BUILD_DIR"
print_info "运行 './build/login' 开始使用"