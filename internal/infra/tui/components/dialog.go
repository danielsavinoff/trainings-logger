package components

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// DialogState describes the action requested by the user.
type DialogState int

const (
	DialogOpen DialogState = iota
	DialogSubmitted
	DialogClosed
)

// Dialog is a reusable modal form with normal and insert modes.
type Dialog struct {
	State         DialogState
	keys          dialogKeyMap
	title         string
	fields        []Field
	active        int
	inserting     bool
	err           string
	width, height int
}

func NewDialog(title string, fields []Field) *Dialog {
	return &Dialog{title: title, fields: fields, keys: defaultDialogKeyMap()}
}

func (d *Dialog) Init() tea.Cmd { return d.focus(0) }

// Value returns the current text of a field.
func (d *Dialog) Value(index int) string { return d.fields[index].input.Value() }

// Invalid keeps the form open and focuses the field that needs correction.
func (d *Dialog) Invalid(field int, message string) tea.Cmd {
	d.State = DialogOpen
	d.err = message

	return d.focus(field)
}

func (d *Dialog) focus(index int) tea.Cmd {
	for i := range d.fields {
		d.fields[i].input.Blur()
	}

	d.active = (index + len(d.fields)) % len(d.fields)
	if d.inserting {
		return d.fields[d.active].input.Focus()
	}

	return nil
}

func (d *Dialog) Update(msg tea.Msg) (*Dialog, tea.Cmd) {
	if _, ok := msg.(tea.MouseMsg); ok {
		return d, nil
	}

	if press, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case key.Matches(press, d.keys.Cancel):
			if d.inserting {
				d.inserting = false
				d.fields[d.active].input.Blur()
			} else {
				d.State = DialogClosed
			}

			return d, nil
		case key.Matches(press, d.keys.Submit):
			d.State = DialogSubmitted

			return d, nil
		}

		if !d.inserting {
			switch {
			case key.Matches(press, d.keys.Down):
				return d, d.focus(d.active + 1)
			case key.Matches(press, d.keys.Up):
				return d, d.focus(d.active - 1)
			case key.Matches(press, d.keys.Insert):
				d.inserting = true

				return d, d.focus(d.active)
			}

			return d, nil
		}

		d.err = ""
	}

	if !d.inserting {
		return d, nil
	}

	field := &d.fields[d.active]
	var cmd tea.Cmd
	field.input, cmd = field.input.Update(msg)

	return d, cmd
}
