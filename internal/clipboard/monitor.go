package clipboard

import (
	"sync"
	"time"
	"unicode/utf8"

	"github.com/atotto/clipboard"
)

// ClipItem represents a clipboard item
type ClipItem struct {
	Content   string
	Timestamp time.Time
}

// Monitor is responsible for monitoring clipboard changes
type Monitor struct {
	history     []ClipItem
	lastContent string
	mutex       sync.RWMutex
	maxItems    int
}

// NewMonitor creates a new clipboard monitor
func NewMonitor() *Monitor {
	return &Monitor{
		history:  make([]ClipItem, 0),
		maxItems: 50, // Store up to 50 clipboard items
	}
}

// Start begins monitoring the clipboard for changes
func (m *Monitor) Start() {
	for {
		content, err := clipboard.ReadAll()
		// 确保内容是有效的UTF-8
		if err == nil && content != "" && content != m.lastContent {
			// 确保内容是有效的UTF-8字符串
			validContent := validateUTF8(content)
			m.addItem(validContent)
			m.lastContent = validContent
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// validateUTF8 确保字符串是有效的UTF-8
func validateUTF8(s string) string {
	if !utf8.ValidString(s) {
		// 如果不是有效的UTF-8字符串，尝试修复
		v := make([]rune, 0, len(s))
		for _, r := range s {
			if r == utf8.RuneError {
				// 跳过无效的UTF-8序列
				continue
			}
			v = append(v, r)
		}
		return string(v)
	}
	return s
}

// addItem adds a new item to the clipboard history
func (m *Monitor) addItem(content string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查是否已存在相同内容
	for i, item := range m.history {
		if item.Content == content {
			// 如果找到相同内容，将其移到队首并更新时间戳
			existingItem := m.history[i]
			existingItem.Timestamp = time.Now()

			// 从原位置删除
			m.history = append(m.history[:i], m.history[i+1:]...)

			// 移到队首
			m.history = append([]ClipItem{existingItem}, m.history...)
			return
		}
	}

	// 如果没有找到相同内容，添加新项目到队首
	m.history = append([]ClipItem{{
		Content:   content,
		Timestamp: time.Now(),
	}}, m.history...)

	// 如果超过最大数量，删除最旧的项目
	if len(m.history) > m.maxItems {
		m.history = m.history[:m.maxItems]
	}
}

// GetHistory returns a copy of the clipboard history
func (m *Monitor) GetHistory() []ClipItem {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Create a copy of the history to avoid race conditions
	historyCopy := make([]ClipItem, len(m.history))
	copy(historyCopy, m.history)

	return historyCopy
}

// CopyToClipboard copies the content to the clipboard
func (m *Monitor) CopyToClipboard(content string) error {
	return clipboard.WriteAll(content)
}
