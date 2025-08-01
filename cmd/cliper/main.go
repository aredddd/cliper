package main

import (
	"fmt"
	"time"

	"github.com/lilithgames/cliper/internal/clipboard"
	"github.com/lilithgames/cliper/internal/ui"
)

func main() {
	// Initialize clipboard monitor
	clipboardMonitor := clipboard.NewMonitor()

	// Start monitoring clipboard
	go clipboardMonitor.Start()

	// 添加恢复机制，防止整个应用崩溃
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("应用发生错误:", r)
			// 给系统一些时间处理后重启应用
			time.Sleep(2 * time.Second)
			main() // 尝试重启
		}
	}()

	// Initialize and run UI
	app := ui.NewApp(clipboardMonitor)
	app.Run()
}
