package views

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m *RecordsModel) focusPanel(next panel) {
	m.activePanel = next
	m.table.Blur()
	if next == exercisesPanel {
		m.table.Focus()
	}
}

func (m RecordsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = size.Width, size.Height
		m = m.resize()

		return m, nil
	}

	if m.dialog != nil {
		return m.updateDialog(msg)
	}

	if wheel, ok := msg.(tea.MouseWheelMsg); ok {
		return m.scrollPanel(wheel)
	}

	// Mouse movement and clicks must not scroll the table's internal viewport.
	if _, ok := msg.(tea.MouseMsg); ok {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.goPending {
			m.goPending = false

			switch msg.String() {
			case "e":
				m.focusPanel(exercisesPanel)

				return m, nil
			case "s":
				m.focusPanel(setsPanel)

				return m, nil
			}
		} else if msg.String() == "g" {
			m.goPending = true

			return m, nil
		}

		switch {
		case msg.String() == "esc", msg.String() == "tab", msg.String() == "shift+tab":
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.New):
			cmd = m.openDialog(false)

			return m, cmd
		case key.Matches(msg, m.keys.Edit):
			cmd = m.openDialog(true)

			return m, cmd
		case key.Matches(msg, m.keys.Delete):
			m.deleteRecord()

			return m, nil
		case msg.String() == "enter":
			if m.activePanel == exercisesPanel {
				m.focusPanel(setsPanel)
			}

			return m, nil
		}
	}

	if press, ok := msg.(tea.KeyPressMsg); ok && m.paginate(press) {
		return m, nil
	}

	if m.activePanel == setsPanel {
		if press, ok := msg.(tea.KeyPressMsg); ok && len(m.setLines) > 0 {
			keys := m.table.KeyMap
			next := m.setCursor

			switch {
			case key.Matches(press, keys.LineUp):
				next--
			case key.Matches(press, keys.LineDown):
				next++
			case key.Matches(press, keys.HalfPageUp):
				next -= max(1, m.timeline.Height()/2)
			case key.Matches(press, keys.HalfPageDown):
				next += max(1, m.timeline.Height()/2)
			case key.Matches(press, keys.GotoTop):
				next = 0
			case key.Matches(press, keys.GotoBottom):
				next = len(m.setLines) - 1
			}

			m.setCursor = min(max(0, next), len(m.setLines)-1)
			m.refreshTimeline()

			return m, nil
		}

		return m, nil
	}

	previous := m.table.Cursor()
	m.table, cmd = m.table.Update(msg)
	m.syncExercisePager()
	if m.table.Cursor() != previous {
		m.refreshSelectedSets()
	}

	return m, cmd
}
