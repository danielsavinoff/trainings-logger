// Package testutils provides input helpers for TUI tests.
package testutils

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Key builds a key press, including control keys and text input.
func Key(key string) tea.KeyPressMsg {
	if strings.HasPrefix(key, "ctrl+") {
		return tea.KeyPressMsg{Code: []rune(strings.TrimPrefix(key, "ctrl+"))[0], Mod: tea.ModCtrl}
	}

	switch key {
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	case "tab":
		return tea.KeyPressMsg{Code: tea.KeyTab}
	case "shift+tab":
		return tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}
	default:
		return tea.KeyPressMsg{Code: []rune(key)[0], Text: key}
	}
}
