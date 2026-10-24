package components

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/paginator"
)

func NewPaginator() paginator.Model {
	p := paginator.New()
	p.KeyMap.PrevPage = key.NewBinding(key.WithKeys(
		"left",
		"h",
		"pgup",
		"b",
	))
	p.KeyMap.NextPage = key.NewBinding(key.WithKeys(
		"right",
		"pgdown",
		"f",
		"space",
	))

	return p
}
