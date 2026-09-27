package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) dialogLayers() ([]*lipgloss.Layer, *tea.Cursor) {
	form := m.dialogView()
	formWidth, formHeight := lipgloss.Size(form)
	x, y := max(0, (m.width-formWidth)/2), max(0, (m.height-formHeight)/2)
	layers := []*lipgloss.Layer{lipgloss.NewLayer(form).ID("dialog").Z(1).X(x).Y(y)}
	return layers, m.dialogCursor(x, y)
}

func (m model) dialogCursor(x, y int) *tea.Cursor {
	width, compact := m.dialogSize()
	fieldWidth := width - 6 - 17
	fieldX, fieldY := 3+17, 2+m.dialog.active*2
	if compact {
		fieldWidth = width - 6
		fieldX, fieldY = 3, 2
	}
	_, caret := recordInputView(m.dialog.fields[m.dialog.active], fieldWidth)
	cursor := tea.NewCursor(x+fieldX+caret, y+fieldY)
	cursor.Color = activePanelColor
	cursor.Shape = tea.CursorBlock
	cursor.Blink = false
	if m.dialog.inserting {
		cursor.Shape = tea.CursorBar
		cursor.Blink = true
	}
	if cursor.Position.X >= m.width || cursor.Position.Y >= m.height {
		return nil
	}
	return cursor
}
