package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caseymrm/menuet"
	"github.com/lilithgames/cliper/internal/clipboard"
)

// App represents the UI application
type App struct {
	monitor *clipboard.Monitor
}

type popupItem struct {
	Text      string `json:"text,omitempty"`
	Index     int    `json:"index,omitempty"`
	Enabled   bool   `json:"enabled,omitempty"`
	Separator bool   `json:"separator,omitempty"`
}

type historyMenuItem struct {
	Content string
	Text    string
}

var activeApp *App

// NewApp creates a new UI application
func NewApp(monitor *clipboard.Monitor) *App {
	return &App{
		monitor: monitor,
	}
}

// Run starts the UI application
func (a *App) Run() {
	activeApp = a
	menuet.App().Name = "Cliper"
	menuet.App().Label = "CL" // 必须设置Label属性，这是应用程序在状态栏显示的标识
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "📎",
	})
	menuet.App().Children = a.menuItems
	startHotkey(a)
	// Disable auto-update to prevent JSON parsing errors
	// menuet.App().AutoUpdate.Version = "1.0.0"
	// menuet.App().AutoUpdate.Repo = "lilithgames/cliper"

	// Setup timed refresh of menu items (every second)
	// 使用更稳定的刷新机制，防止应用从状态栏消失
	go func() {
		for {
			menuet.App().SetMenuState(&menuet.MenuState{
				Title: "📎",
			})
			menuet.App().MenuChanged()
			time.Sleep(1 * time.Second)
		}
	}()

	// Start the app
	menuet.App().RunApplication()
}

// menuItems returns the menu items for the status bar menu
func (a *App) menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}

	// Add header
	items = append(items, menuet.MenuItem{
		Text: "Cliper - Clipboard History",
	})

	// Add separator
	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})

	// Get clipboard history
	history := a.monitor.GetHistory()

	// If history is empty, show a message
	if len(history) == 0 {
		items = append(items, menuet.MenuItem{
			Text: "No clipboard history yet",
		})
	} else {
		for _, item := range buildHistoryMenuItems(history) {
			items = append(items, menuet.MenuItem{
				Text:    item.Text,
				Clicked: a.createClickHandler(item.Content),
			})
		}
	}

	// Add separator
	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})

	// 添加关于选项
	items = append(items, menuet.MenuItem{
		Text: "关于Cliper",
		Clicked: func() {
			// 显示关于对话框
			menuet.App().Alert(menuet.Alert{
				MessageText:     "关于Cliper",
				InformativeText: "Cliper - 轻量级剪贴板历史工具\n\nhttps://github.com/aredddd/cliper",
			})
		},
	})

	// 分隔线
	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})

	// Add quit button
	items = append(items, menuet.MenuItem{
		Text: "退出",
		Clicked: func() {
			// Use standard Go way to exit the application
			os.Exit(0)
		},
	})

	return items
}

// createClickHandler creates a click handler for a clipboard item
func (a *App) createClickHandler(content string) func() {
	return func() {
		a.monitor.CopyToClipboard(content)
	}
}

func (a *App) popupItemsJSON() string {
	items := []popupItem{{Text: "Cliper - Clipboard History"}}
	historyItems := buildHistoryMenuItems(a.monitor.GetHistory())

	if len(historyItems) == 0 {
		items = append(items, popupItem{Text: "No clipboard history yet"})
	} else {
		items = append(items, popupItem{Separator: true})
		for i, item := range historyItems {
			items = append(items, popupItem{
				Text:    item.Text,
				Index:   i,
				Enabled: true,
			})
		}
	}

	data, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func (a *App) copyHistoryItem(index int) {
	items := buildHistoryMenuItems(a.monitor.GetHistory())
	if index < 0 || index >= len(items) {
		return
	}
	_ = a.monitor.CopyToClipboard(items[index].Content)
}

func (a *App) pasteHistoryItem(index int) {
	a.copyHistoryItem(index)
	pasteFromHotkeyMenu()
}

func buildHistoryMenuItems(history []clipboard.ClipItem) []historyMenuItem {
	limit := len(history)
	if limit > 20 {
		limit = 20
	}

	items := make([]historyMenuItem, 0, limit)
	for i := 0; i < limit; i++ {
		item := history[i]
		menuText := fmt.Sprintf("%s (%s)", formatDisplayText(item.Content), formatTimeAgo(item.Timestamp))
		items = append(items, historyMenuItem{
			Content: item.Content,
			Text:    menuText,
		})
	}
	return items
}

func formatDisplayText(content string) string {
	displayText := strings.ReplaceAll(content, "\n", " ")

	runeCount := 0
	for i := range displayText {
		runeCount++
		if runeCount > 40 {
			return displayText[:i] + "..."
		}
	}
	return displayText
}

// formatTimeAgo formats a timestamp as a human-readable time ago string
func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	default:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
