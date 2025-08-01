package ui

import (
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

// NewApp creates a new UI application
func NewApp(monitor *clipboard.Monitor) *App {
	return &App{
		monitor: monitor,
	}
}

// Run starts the UI application
func (a *App) Run() {
	menuet.App().Name = "Cliper"
	menuet.App().Label = "📋"
	menuet.App().Children = a.menuItems
	menuet.App().AutoUpdate.Version = "1.0.0"
	menuet.App().AutoUpdate.Repo = "lilithgames/cliper"

	// Setup timed refresh of menu items (every second)
	go func() {
		for {
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
		// Add history items
		for i, item := range history {
			// Limit to 20 items in the menu for performance
			if i >= 20 {
				break
			}

			// Truncate content if too long
			displayText := item.Content
			if len(displayText) > 60 {
				displayText = displayText[:57] + "..."
			}

			// Replace newlines with spaces for display
			displayText = strings.ReplaceAll(displayText, "\n", " ")

			// Format timestamp
			timeAgo := formatTimeAgo(item.Timestamp)

			// Create menu item with timestamp in the menu text
			menuText := fmt.Sprintf("%s (%s)", displayText, timeAgo)
			items = append(items, menuet.MenuItem{
				Text:    menuText,
				Clicked: a.createClickHandler(item.Content),
			})
		}
	}

	// Add separator
	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})

	// Add quit button
	items = append(items, menuet.MenuItem{
		Text: "Quit",
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
