// START_MODULE_CONTRACT
//   PURPOSE: Route records, dialog, help, questionnaire, and asynchronous transfers.
//   SCOPE: Preserve input ownership; 1/2 focus panels, r reopens setup, s saves, and l loads.
//   DEPENDS: M-TUI-HELP, M-TUI-QUESTIONNAIRE, M-APP-REMOTE-STATE
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   focusPanel - change active records panel
//   Update - dispatch events to the active interaction mode
// END_MODULE_MAP

package views

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// START_BLOCK_RECORDS_INPUT
func (m *RecordsModel) focusPanel(next panel) {
	m.activePanel = next
	m.table.Blur()
	if next == exercisesPanel {
		m.table.Focus()
	}
}

func (m RecordsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if result, ok := msg.(recordTransferMsg); ok {
		return m.finishTransfer(result), nil
	}

	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = size.Width, size.Height
		m = m.resize()

		return m, nil
	}

	if m.dialog != nil {
		return m.updateDialog(msg)
	}

	if m.showingHelp {
		if _, ok := msg.(tea.KeyPressMsg); ok {
			m.showingHelp = false
			return m, nil
		}
		m.help = m.help.update(msg)
		return m, nil
	}
	if m.questionnaire != nil {
		return m.updateQuestionnaire(msg)
	}

	if press, ok := msg.(tea.KeyPressMsg); ok && key.Matches(
		press,
		m.keys.Help,
	) {
		m.goPending = false
		m.help = newHelpModel(
			m.width,
			m.height,
		)
		m.showingHelp = true
		return m, nil
	}
	if m.transfer != nil && m.transfer.load {
		if press, ok := msg.(tea.KeyPressMsg); ok && key.Matches(
			press,
			m.keys.Quit,
		) {
			return m, tea.Quit
		}
		return m, nil
	}

	if _, ok := msg.(tea.MouseMsg); ok {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.transfer != nil {
			switch msg.String() {
			case "n", "e", "d", "r":
				return m, nil
			}
		}
		switch msg.String() {
		case "1":
			m.goPending = false
			m.focusPanel(exercisesPanel)

			return m, nil
		case "2":
			m.goPending = false
			m.focusPanel(setsPanel)

			return m, nil
		}

		if m.goPending {
			m.goPending = false

			switch msg.String() {
			case "g":
				if m.activePanel == exercisesPanel {
					m.table.GotoTop()
					m.syncExercisePager()
					m.refreshSelectedSets()
				} else if len(m.setLines) > 0 {
					m.setCursor = 0
					m.refreshTimeline()
				}

				return m, nil
			}
		} else if msg.String() == "g" {
			m.goPending = true

			return m, nil
		}
		if msg.String() == "r" {
			cmd := m.openQuestionnaire()

			return m, cmd
		}
		if msg.String() == "s" || msg.String() == "l" {
			cmd := m.startTransfer(msg.String() == "l")
			return m, cmd
		}

		switch {
		case msg.String() == "esc", msg.String() == "tab", msg.String() == "shift+tab":
			return m, nil
		case key.Matches(
			msg,
			m.keys.Quit,
		):
			return m, tea.Quit
		case key.Matches(
			msg,
			m.keys.New,
		):
			cmd = m.openDialog(false)

			return m, cmd
		case key.Matches(
			msg,
			m.keys.Edit,
		):
			cmd = m.openDialog(true)

			return m, cmd
		case key.Matches(
			msg,
			m.keys.Delete,
		):
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
			case key.Matches(
				press,
				keys.LineUp,
			):
				next--
			case key.Matches(
				press,
				keys.LineDown,
			):
				next++
			case key.Matches(
				press,
				keys.HalfPageUp,
			):
				next -= max(
					1,
					m.timeline.Height()/2,
				)
			case key.Matches(
				press,
				keys.HalfPageDown,
			):
				next += max(
					1,
					m.timeline.Height()/2,
				)
			case key.Matches(
				press,
				keys.GotoTop,
			):
				next = 0
			case key.Matches(
				press,
				keys.GotoBottom,
			):
				next = len(m.setLines) - 1
			}

			m.setCursor = min(
				max(
					0,
					next,
				),
				len(m.setLines)-1,
			)
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

// END_BLOCK_RECORDS_INPUT

func (m RecordsModel) HelpActive() bool {
	return m.showingHelp
}
