//go:build darwin

package clipboard

import (
	"os/exec"
)

var readClipboardText = readMacClipboardText

func readMacClipboardText() (string, error) {
	out, err := exec.Command("pbpaste", "-Prefer", "txt").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
