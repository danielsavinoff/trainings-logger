package main

import (
	"image"

	tea "charm.land/bubbletea/v2"
)

// Item bounds exclude borders, padding, column headings, pagination and help.
func (m model) itemBounds(p panel) image.Rectangle {
	x := panelSpacing.GetPaddingLeft() + baseStyle.GetBorderLeftSize() + baseStyle.GetPaddingLeft()
	y := panelSpacing.GetPaddingTop() + baseStyle.GetBorderTopSize() + baseStyle.GetPaddingTop() + 1
	width, height := m.table.Width(), m.table.Height()
	if p == setsPanel {
		x += m.width / 2
		width, height = m.timeline.Width(), m.timeline.Height()
	}
	return image.Rect(x, y, x+width, y+height)
}

func (m model) scrollPanel(msg tea.MouseWheelMsg) (tea.Model, tea.Cmd) {
	point := image.Pt(msg.X, msg.Y)
	for _, p := range []panel{exercisesPanel, setsPanel} {
		if !point.In(m.itemBounds(p)) {
			continue
		}
		var code rune
		switch msg.Button {
		case tea.MouseWheelUp:
			code = tea.KeyUp
		case tea.MouseWheelDown:
			code = tea.KeyDown
		case tea.MouseWheelLeft:
			code = tea.KeyLeft
		case tea.MouseWheelRight:
			code = tea.KeyRight
		default:
			return m, nil
		}
		m.goPending = false
		m.focusPanel(p)
		return m.Update(tea.KeyPressMsg{Code: code})
	}
	return m, nil
}
