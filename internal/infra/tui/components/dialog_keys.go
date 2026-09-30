package components

import "charm.land/bubbles/v2/key"

type dialogKeyMap struct {
	Up, Down, Cancel, Insert, Submit key.Binding
}

func defaultDialogKeyMap() dialogKeyMap {
	return dialogKeyMap{
		Up:     key.NewBinding(key.WithKeys("k", "up", "shift+tab"), key.WithHelp("↑/k", "move up")),
		Down:   key.NewBinding(key.WithKeys("j", "down", "tab"), key.WithHelp("↓/j", "move down")),
		Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "normal mode")),
		Insert: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "insert mode")),
		Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
	}
}
