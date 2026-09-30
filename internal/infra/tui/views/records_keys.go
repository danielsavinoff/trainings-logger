package views

import "charm.land/bubbles/v2/key"

type recordsKeyMap struct {
	New, Edit, Delete, Quit key.Binding
}

func defaultRecordsKeyMap() recordsKeyMap {
	return recordsKeyMap{
		New:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new record")),
		Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit record")),
		Delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete record")),
		Quit:   key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}
}
