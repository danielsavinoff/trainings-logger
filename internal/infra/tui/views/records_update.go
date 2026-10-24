// START_MODULE_CONTRACT
//   PURPOSE: Route records, catalog, dialog, help, questionnaire, and asynchronous transfers.
//   SCOPE: Preserve input ownership; 1/2/3 focus records panels, r reopens setup, s saves, and l loads.
//   DEPENDS: M-TUI-CATALOG, M-TUI-HELP, M-TUI-QUESTIONNAIRE, M-APP-REMOTE-STATE
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
	m.catalogTable.Blur()
	m.performanceTable.Blur()
	m.sessionTable.Blur()
	m.setTable.Blur()
	m.activeTable().Focus()
}

func (m RecordsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if result, ok := msg.(recordTransferMsg); ok {
		return m.finishTransfer(result), nil
	}

	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = max(
			0,
			size.Width,
		), max(
			0,
			size.Height,
		)
		return m.resize(), nil
	}

	if m.dialog != nil {
		return m.updateDialog(msg)
	}

	if m.showingCatalog {
		return m.updateCatalog(msg)
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

	if press, ok := msg.(tea.KeyPressMsg); ok {
		if m.transfer != nil {
			switch press.String() {
			case "n", "e", "d", "r":
				return m, nil
			}
		}
		switch press.String() {
		case "1":
			m.goPending = false
			m.focusPanel(sessionsPanel)
			return m, nil
		case "2":
			m.goPending = false
			m.focusPanel(performancesPanel)
			return m, nil
		case "3":
			m.goPending = false
			m.focusPanel(setsPanel)
			return m, nil
		}

		if m.goPending {
			m.goPending = false
			if press.String() == "g" {
				m.activeTable().GotoTop()
				m.selectionChanged()
				return m, nil
			}
		} else if press.String() == "g" {
			m.goPending = true
			return m, nil
		}
		switch {
		case press.String() == "r":
			cmd := m.openQuestionnaire()
			return m, cmd
		case press.String() == "s" || press.String() == "l":
			cmd := m.startTransfer(press.String() == "l")
			return m, cmd
		case press.String() == "esc", press.String() == "tab", press.String() == "shift+tab":
			return m, nil
		case key.Matches(
			press,
			m.keys.Quit,
		):
			return m, tea.Quit
		case key.Matches(
			press,
			m.keys.New,
		):
			cmd := m.openDialog(false)
			return m, cmd
		case key.Matches(
			press,
			m.keys.Edit,
		):
			cmd := m.openDialog(true)
			return m, cmd
		case key.Matches(
			press,
			m.keys.Delete,
		):
			m.deleteRecord()
			return m, nil
		case press.String() == "enter":
			switch m.activePanel {
			case sessionsPanel:
				m.focusPanel(performancesPanel)
			case performancesPanel:
				m.focusPanel(setsPanel)
			}

			return m, nil
		}

		if m.paginate(press) {
			return m, nil
		}
	}

	t := m.activeTable()
	previous := t.Cursor()
	var cmd tea.Cmd
	*t, cmd = t.Update(msg)
	m.syncPagers()
	if t.Cursor() != previous {
		m.selectionChanged()
	}

	return m, cmd
}

// END_BLOCK_RECORDS_INPUT
func (m RecordsModel) HelpActive() bool { return m.showingHelp }
