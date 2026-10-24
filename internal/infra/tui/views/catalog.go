// START_MODULE_CONTRACT
//   PURPOSE: Select and manage catalog exercises in an isolated full-screen view.
//   SCOPE: Catalog performance/return, local errors, keyboard navigation, selection, and CRUD routing.
//   DEPENDS: M-TUI-RECORDS, M-APP-PERFORMANCES, M-DOMAIN-EXERCISE
//   LINKS: M-TUI-CATALOG, V-M-TUI-CATALOG
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   openCatalog - enter catalog while retaining records errors and selections
//   closeCatalog - return to workout exercises and restore records error state
//   updateCatalog - isolate catalog input, selection, and definition management
//   catalogStatusView - render only actual catalog errors with table spacing
//   catalogView - compose the full-screen table and optional local errors
// END_MODULE_MAP

package views

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// START_BLOCK_CATALOG_VIEW
func (m *RecordsModel) openCatalog() {
	if m.showingCatalog || m.transfer != nil {
		return
	}
	sessionID, catalogID, performanceID, setID := m.selectionIDs()
	m.catalogRecordsError = m.err
	m.err = ""
	m.showingCatalog = true
	m.goPending = false
	m.focusPanel(catalogPanel)
	if err := m.reloadRecords(
		context.Background(),
		sessionID,
		catalogID,
		performanceID,
		setID,
	); err != nil {
		m.err = err.Error()
	}
	*m = m.resize()
}

func (m *RecordsModel) closeCatalog() {
	m.showingCatalog = false
	m.err = m.catalogRecordsError
	m.catalogRecordsError = ""
	m.goPending = false
	m.focusPanel(performancesPanel)
	*m = m.resize()
}

func (m RecordsModel) updateCatalog(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.MouseMsg); ok {
		return m, nil
	}

	if press, ok := msg.(tea.KeyPressMsg); ok {
		switch press.String() {
		case "esc":
			m.closeCatalog()
			return m, nil
		case "enter":
			m.goPending = false
			m.addCatalogExercise()
			return m, nil
		case "n":
			cmd := m.openDialog(false)
			return m, cmd
		case "e":
			cmd := m.openDialog(true)
			return m, cmd
		case "d":
			m.goPending = false
			m.deleteRecord()
			return m, nil
		}

		if m.goPending {
			m.goPending = false
			if press.String() == "g" {
				m.catalogTable.GotoTop()
				m.syncPagers()
				return m, nil
			}
		} else if press.String() == "g" {
			m.goPending = true
			return m, nil
		}

		if m.paginate(press) {
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.catalogTable, cmd = m.catalogTable.Update(msg)
	m.syncPagers()
	return m, cmd
}

func (m RecordsModel) catalogStatusView() string {
	if m.err == "" {
		return ""
	}
	return recordsPaneSpacing.Render(m.footerLine(m.err))
}

func (m RecordsModel) catalogView() string {
	pane := m.recordsPane(
		catalogPanel,
		"Exercise catalog",
		m.catalogTable,
		m.catalogPager,
	)
	status := m.catalogStatusView()
	if status == "" {
		return pane
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		pane,
		status,
	)
}

// END_BLOCK_CATALOG_VIEW
