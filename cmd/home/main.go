package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
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
	
	// 6. 保持程序运行，等待用户输入或浏览器关闭，保持程序运行状态
	waitForUserInputOrBrowserClose(page)
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

// takeScreenshot 截取指定元素的截图
func takeScreenshot(page *rod.Page, selector string) {
	logrus.Infof("正在尝试为元素 '%s' 截图...", selector)

	// 检查元素是否存在
	if has, _, err := page.Has(selector); !has || err != nil {
		logrus.Errorf("找不到要截图的元素 '%s'，错误: %v", selector, err)
		return
	}

	// 获取元素并截图
	element := page.MustElement(selector)
	screenshot, err := element.Screenshot(proto.PageCaptureScreenshotFormatPng, 0)
	if err != nil {
		logrus.Errorf("为元素 '%s' 截图失败: %v", selector, err)
		return
	}

	// 将截图保存到文件
	shotDir := "shot"
	if err := os.MkdirAll(shotDir, 0755); err != nil {
		logrus.Errorf("创建截图目录失败: %v", err)
		return
	}
	filePath := fmt.Sprintf("shot/shot_%s.png", time.Now().Format("20060102150405"))
	if err := os.WriteFile(filePath, screenshot, 0644); err != nil {
		logrus.Errorf("保存截图失败: %v", err)
		return
	}

	logrus.Infof("截图成功！已保存到: %s", filePath)
}

// isBrowserClosed 检查浏览器是否已关闭
func isBrowserClosed(page *rod.Page) bool {
	defer func() {
		if r := recover(); r != nil {
			// 如果发生 panic，说明浏览器已关闭
			logrus.Debug("检测到浏览器关闭（通过 panic 恢复）")
		}
	}()
	
	// 尝试获取页面信息，如果浏览器关闭会抛出异常
	_, err := page.Info()
	return err != nil
}

// monitorBrowserStatus 监控浏览器状态
func monitorBrowserStatus(page *rod.Page, done chan<- bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if isBrowserClosed(page) {
				logrus.Info("检测到浏览器已关闭，程序即将退出...")
				done <- true
				return
			}
		}
	}
}

// waitForUserInputOrBrowserClose 等待用户输入命令或浏览器关闭，保持程序运行状态
// 提供交互式命令行界面，允许用户查看帮助或退出程序
// 同时监控浏览器状态，如果浏览器关闭则自动退出程序
func waitForUserInputOrBrowserClose(page *rod.Page) {
	fmt.Println("\n=== 小红书首页交互程序 ===")
	fmt.Println("这是一个独立的首页交互程序")
	fmt.Println("输入 'quit' 或 'exit' 退出程序，或按 Ctrl+C 强制退出")
	fmt.Println("关闭浏览器窗口也会自动退出程序")
	fmt.Print("请输入命令: ")
	
	// 创建通道用于协调不同的监听器
	done := make(chan bool, 1)
	var wg sync.WaitGroup
	
	// 启动浏览器状态监控 goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitorBrowserStatus(page, done)
	}()
	
	// 启动用户输入监听 goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := strings.TrimSpace(strings.ToLower(scanner.Text()))
			
			switch input {
			case "quit", "exit", "q":
				fmt.Println("正在退出程序...")
				done <- true
				return
			case "help", "h":
				fmt.Println("可用命令:")
				fmt.Println("  quit/exit/q - 退出程序")
				fmt.Println("  help/h - 显示帮助信息")
				fmt.Println("  status - 显示程序状态")
				fmt.Println("  browser - 检查浏览器状态")
				fmt.Println("  shot - 截取用户信息的截图")
			case "status":
				fmt.Println("程序状态: 运行中")
				fmt.Println("版本: 1.0.0")
				fmt.Println("功能: 小红书首页交互")
			case "browser":
				if isBrowserClosed(page) {
					fmt.Println("浏览器状态: 已关闭")
				} else {
					fmt.Println("浏览器状态: 运行中")
				}
			case "shot":
				takeScreenshot(page, ".user-info")
			case "":
				// 空输入，继续等待
			default:
				fmt.Printf("未知命令: %s，输入 'help' 查看可用命令\n", input)
			}
			
			// 检查是否需要退出
			select {
			case <-done:
				return
			default:
				fmt.Print("请输入命令: ")
			}
		}
		
		if err := scanner.Err(); err != nil {
			logrus.Errorf("读取输入时出错: %v", err)
		}
	}()
	
	// 等待任一监听器发出退出信号
	<-done
	
	// 通知所有 goroutine 退出
	close(done)
	
	// 等待所有 goroutine 完成
	wg.Wait()
	
	logrus.Info("程序已退出")
}