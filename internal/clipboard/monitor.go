package clipboard

import (
	"sync"
	"time"

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
		if err == nil && content != "" && content != m.lastContent {
			m.addItem(content)
			m.lastContent = content
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// addItem adds a new item to the clipboard history
func (m *Monitor) addItem(content string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Add new item to the beginning of the slice
	m.history = append([]ClipItem{{
		Content:   content,
		Timestamp: time.Now(),
	}}, m.history...)

	// Trim history if it exceeds max items
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
