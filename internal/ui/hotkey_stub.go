//go:build !darwin

package ui

func startHotkey(app *App) {
	activeApp = app
}
