package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
	"github.com/xpzouying/xiaohongshu-mcp/cookies"
	"github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"
)

func main() {
	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	
	// 启动首页程序
	handleHomePage()
}

// handleHomePage 处理已登录用户的首页逻辑
// 读取 cookies.json 文件，打开浏览器，导航到小红书首页和个人页面
func handleHomePage() {
	logrus.Info("欢迎使用小红书首页交互程序")
	
	// 1. 初始化浏览器
	logrus.Info("正在初始化浏览器...")
	browserInstance := browser.NewBrowser(false) // false 表示非无头模式，可以看到浏览器界面
	defer browserInstance.Close()
	
	page := browserInstance.NewPage()
	defer page.Close()
	
	// 2. 检查并加载 cookies
	logrus.Info("正在检查 cookies 文件...")
	cookiesPath := cookies.GetCookiesFilePath()
	if _, err := os.Stat(cookiesPath); os.IsNotExist(err) {
		logrus.Fatalf("cookies 文件不存在: %s，请先运行登录程序", cookiesPath)
	}
	
	logrus.Infof("找到 cookies 文件: %s", cookiesPath)
	
	// 3. 导航到小红书首页
	logrus.Info("正在导航到小红书首页...")
	navigate := xiaohongshu.NewNavigate(page)
	if err := navigate.ToExplorePage(context.Background()); err != nil {
		logrus.Fatalf("导航到首页失败: %v", err)
	}
	
	// 等待页面加载
	time.Sleep(2 * time.Second)
	
	// 4. 检查登录状态
	logrus.Info("正在检查登录状态...")
	loginAction := xiaohongshu.NewLogin(page)
	isLoggedIn, err := loginAction.CheckLoginStatus(context.Background())
	if err != nil {
		logrus.Fatalf("检查登录状态失败: %v", err)
	}
	
	if !isLoggedIn {
		logrus.Fatal("未检测到登录状态，请先运行登录程序")
	}
	
	logrus.Info("登录状态确认成功！")
	
	// 5. 导航到个人页面（"我的链接"）
	logrus.Info("正在导航到个人页面...")
	if err := navigateToMyProfile(page); err != nil {
		logrus.Errorf("导航到个人页面失败: %v", err)
	} else {
		logrus.Info("成功导航到个人页面！")
	}
	
	// 6. 保持程序运行，等待用户输入
	waitForUserInput()
}

// navigateToMyProfile 导航到个人页面
func navigateToMyProfile(page *rod.Page) error {
	logrus.Info("正在导航到个人页面...")
	
	// 等待页面加载完成
	time.Sleep(2 * time.Second)
	
	// 方法1: 尝试点击右上角的用户头像
	userAvatarSelectors := []string{
		".main-container .user .link-wrapper .channel",
		".user-info .avatar",
		".header .user-avatar",
		"[data-testid='user-avatar']",
		".user .avatar",
	}
	
	for _, selector := range userAvatarSelectors {
		if exists, _, _ := page.Has(selector); exists {
			logrus.Infof("找到用户头像元素: %s，正在点击...", selector)
			page.MustElement(selector).MustClick()
			time.Sleep(3 * time.Second)
			
			// 检查是否成功进入个人页面
			currentURL := page.MustInfo().URL
			if currentURL != "https://www.xiaohongshu.com/explore" {
				logrus.Infof("成功进入个人页面: %s", currentURL)
				return nil
			}
		}
	}
	
	// 方法2: 尝试通过菜单导航
	logrus.Info("尝试通过菜单导航到个人页面...")
	menuSelectors := []string{
		".nav-menu .profile",
		".sidebar .my-profile",
		"[href*='/user/profile']",
	}
	
	for _, selector := range menuSelectors {
		if exists, _, _ := page.Has(selector); exists {
			logrus.Infof("找到菜单元素: %s，正在点击...", selector)
			page.MustElement(selector).MustClick()
			time.Sleep(3 * time.Second)
			return nil
		}
	}
	
	// 方法3: 尝试直接导航到个人页面URL（通用路径）
	logrus.Info("尝试直接导航到个人页面...")
	page.MustNavigate("https://www.xiaohongshu.com/user/profile").MustWaitLoad()
	time.Sleep(2 * time.Second)
	
	logrus.Info("已尝试导航到个人页面")
	return nil
}

// waitForUserInput 等待用户输入命令，保持程序运行状态
// 提供交互式命令行界面，允许用户查看帮助或退出程序
func waitForUserInput() {
	fmt.Println("\n=== 小红书首页交互程序 ===")
	fmt.Println("这是一个独立的首页交互程序")
	fmt.Println("输入 'quit' 或 'exit' 退出程序，或按 Ctrl+C 强制退出")
	fmt.Print("请输入命令: ")
	
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		
		switch input {
		case "quit", "exit", "q":
			fmt.Println("正在退出程序...")
			return
		case "help", "h":
			fmt.Println("可用命令:")
			fmt.Println("  quit/exit/q - 退出程序")
			fmt.Println("  help/h - 显示帮助信息")
			fmt.Println("  status - 显示程序状态")
		case "status":
			fmt.Println("程序状态: 运行中")
			fmt.Println("版本: 1.0.0")
			fmt.Println("功能: 小红书首页交互")
		case "":
			// 空输入，继续等待
		default:
			fmt.Printf("未知命令: %s，输入 'help' 查看可用命令\n", input)
		}
		
		fmt.Print("请输入命令: ")
	}
	
	if err := scanner.Err(); err != nil {
		logrus.Errorf("读取输入时出错: %v", err)
	}
}