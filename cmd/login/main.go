package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
	"github.com/xpzouying/xiaohongshu-mcp/cookies"
	"github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"
)

func main() {
	var (
		binPath string // 浏览器二进制文件路径
	)
	flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
	flag.Parse()

	// 检查cookies.json文件是否存在
	cookiesPath := cookies.GetCookiesFilePath()
	if _, err := os.Stat(cookiesPath); err == nil {
		logrus.Infof("发现cookies文件: %s，尝试使用已保存的登录状态", cookiesPath)
		
		// 登录的时候，需要界面，所以不能无头模式
		b := browser.NewBrowser(false, browser.WithBinPath(binPath))
		defer b.Close()

		page := b.NewPage()
		defer page.Close()

		// 加载cookies
		if err := loadCookies(page); err != nil {
			logrus.Warnf("加载cookies失败: %v，将进行重新登录", err)
		} else {
			// 检查登录状态
			action := xiaohongshu.NewLogin(page)
			status, err := action.CheckLoginStatus(context.Background())
			if err != nil {
				logrus.Warnf("检查登录状态失败: %v，将进行重新登录", err)
			} else if status {
				logrus.Info("使用已保存的cookies登录成功，直接跳转到首页")
				// 导航到首页
				page.MustNavigate("https://www.xiaohongshu.com/explore").MustWaitLoad()
				logrus.Info("已跳转到小红书首页")
				
				// 保持程序运行，等待用户输入
				waitForUserInput()
			} else {
				logrus.Info("已保存的cookies已失效，将进行重新登录")
			}
		}
	} else {
		logrus.Info("未发现cookies文件，将进行登录流程")
	}

	// 登录的时候，需要界面，所以不能无头模式
	b := browser.NewBrowser(false, browser.WithBinPath(binPath))
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	action := xiaohongshu.NewLogin(page)

	status, err := action.CheckLoginStatus(context.Background())
	if err != nil {
		logrus.Fatalf("failed to check login status: %v", err)
	}

	logrus.Infof("当前登录状态: %v", status)

	if status {
		logrus.Info("检测到已登录状态，保持程序运行")
		// 导航到首页
		page.MustNavigate("https://www.xiaohongshu.com/explore").MustWaitLoad()
		logrus.Info("已跳转到小红书首页")
		
		// 保持程序运行，等待用户输入
		waitForUserInput()
		return
	}

	// 开始登录流程
	logrus.Info("开始登录流程...")
	if err = action.Login(context.Background()); err != nil {
		logrus.Fatalf("登录失败: %v", err)
	} else {
		if err := saveCookies(page); err != nil {
			logrus.Fatalf("failed to save cookies: %v", err)
		}
	}

	// 再次检查登录状态确认成功
	status, err = action.CheckLoginStatus(context.Background())
	if err != nil {
		logrus.Fatalf("failed to check login status after login: %v", err)
	}

	if status {
		logrus.Info("登录成功！")
		// 导航到首页
		page.MustNavigate("https://www.xiaohongshu.com/explore").MustWaitLoad()
		logrus.Info("已跳转到小红书首页")
		
		// 保持程序运行，等待用户输入
		waitForUserInput()
	} else {
		logrus.Error("登录流程完成但仍未登录")
	}

}

func saveCookies(page *rod.Page) error {
	cks, err := page.Browser().GetCookies()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cks)
	if err != nil {
		return err
	}

	cookieLoader := cookies.NewLoadCookie(cookies.GetCookiesFilePath())
	return cookieLoader.SaveCookies(data)
}

// loadCookies 从文件加载cookies到浏览器
func loadCookies(page *rod.Page) error {
	cookieLoader := cookies.NewLoadCookie(cookies.GetCookiesFilePath())
	data, err := cookieLoader.LoadCookies()
	if err != nil {
		return err
	}

	var cks []*proto.NetworkCookie
	if err := json.Unmarshal(data, &cks); err != nil {
		return err
	}

	// 转换为 NetworkCookieParam 类型
	var cookieParams []*proto.NetworkCookieParam
	for _, ck := range cks {
		cookieParams = append(cookieParams, &proto.NetworkCookieParam{
			Name:     ck.Name,
			Value:    ck.Value,
			Domain:   ck.Domain,
			Path:     ck.Path,
			Secure:   ck.Secure,
			HTTPOnly: ck.HTTPOnly,
			SameSite: ck.SameSite,
			Expires:  ck.Expires,
		})
	}

	// 设置cookies到浏览器
	return page.Browser().SetCookies(cookieParams)
}

// waitForUserInput 等待用户输入，保持程序运行
func waitForUserInput() {
	fmt.Println("\n=== 小红书登录程序 ===")
	fmt.Println("浏览器将保持打开状态，你可以继续使用小红书")
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
