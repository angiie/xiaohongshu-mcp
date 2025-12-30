# 命令行工具

这个文件夹包含了小红书 MCP 项目的命令行工具。

## 工具说明

### login - 登录工具
位置: `cmd/login/main.go`

**功能:**
- 处理小红书登录流程
- 保存登录状态到 cookies.json
- 检查现有登录状态

**使用方法:**
```bash
go run ./cmd/login
# 或者编译后运行
go build -o login ./cmd/login
./login
```

**参数:**
- `-bin`: 指定浏览器二进制文件路径（可选）

### home - 首页交互工具
位置: `cmd/home/main.go`

**功能:**
- 读取已保存的登录状态
- 打开小红书首页
- 导航到个人页面
- 提供交互式命令行界面

**使用方法:**
```bash
go run ./cmd/home
# 或者编译后运行
go build -o home ./cmd/home
./home
```

**前置条件:**
- 必须先运行 login 工具完成登录
- 需要存在有效的 cookies.json 文件

## 工作流程

1. 首次使用时，运行 `login` 工具完成登录
2. 登录成功后，可以运行 `home` 工具进行首页交互
3. 如果登录状态失效，重新运行 `login` 工具

## 注意事项

- 两个工具都会以非无头模式运行浏览器（可以看到浏览器界面）
- cookies.json 文件会自动保存在项目根目录
- 程序运行时请不要手动关闭浏览器窗口