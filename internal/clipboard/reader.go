//go:build !darwin

package clipboard

import clipboardlib "github.com/atotto/clipboard"

var readClipboardText = clipboardlib.ReadAll
