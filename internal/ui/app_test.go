package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/lilithgames/cliper/internal/clipboard"
)

func TestFormatDisplayText(t *testing.T) {
	t.Parallel()

	got := formatDisplayText("第一行\n第二行")
	if got != "第一行 第二行" {
		t.Fatalf("expected newline to be replaced, got %q", got)
	}

	longText := strings.Repeat("界", 41)
	got = formatDisplayText(longText)
	if got != strings.Repeat("界", 40)+"..." {
		t.Fatalf("expected unicode-safe truncation, got %q", got)
	}
}

func TestBuildHistoryMenuItemsLimitsToTwenty(t *testing.T) {
	t.Parallel()

	history := make([]clipboard.ClipItem, 25)
	now := time.Now()
	for i := range history {
		history[i] = clipboard.ClipItem{
			Content:   "item",
			Timestamp: now,
		}
	}

	items := buildHistoryMenuItems(history)
	if len(items) != 20 {
		t.Fatalf("expected 20 menu items, got %d", len(items))
	}
}
