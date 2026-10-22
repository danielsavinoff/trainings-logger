package components

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (d *Dialog) dialogSize() (width int, compact bool) {
	width = max(12, min(64, d.width-4))

	return width, d.height < 16 || width < 44
}

func (d *Dialog) Resize(width, height int) {
	d.width, d.height = width, height

	width, compact := d.dialogSize()
	inputWidth := width - 6 - 20 // Padding, label, prompt and cursor.
	if compact {
		inputWidth = width - 6 - 3
	}

	for i := range d.fields {
		input := &d.fields[i].input
		position := input.Position()

		input.SetWidth(0)
		input.SetCursor(position)
		input.SetWidth(max(1, inputWidth))
		input.CursorEnd()
		input.SetCursor(position)
	}
}

func (d *Dialog) Layers() ([]*lipgloss.Layer, *tea.Cursor) {
	form := d.View()
	formWidth, formHeight := lipgloss.Size(form)
	x, y := max(0, (d.width-formWidth)/2), max(0, (d.height-formHeight)/2)
	layers := []*lipgloss.Layer{lipgloss.NewLayer(form).ID("dialog").Z(1).X(x).Y(y)}

	return layers, d.dialogCursor(x, y)
}

func (d *Dialog) dialogCursor(x, y int) *tea.Cursor {
	width, compact := d.dialogSize()
	fieldWidth := width - 6 - 17
	fieldX, fieldY := 3+17, 2+d.active*2
	if compact {
		fieldWidth = width - 6
		fieldX, fieldY = 3, 2
	}

	_, caret := fieldView(d.fields[d.active], fieldWidth)
	cursor := tea.NewCursor(x+fieldX+caret, y+fieldY)
	cursor.Color = ActivePanelColor
	cursor.Shape = tea.CursorBlock
	cursor.Blink = false
	if d.inserting {
		cursor.Shape = tea.CursorBar
		cursor.Blink = true
	}

	if cursor.Position.X >= d.width || cursor.Position.Y >= d.height {
		return nil
	}

	return cursor
}
