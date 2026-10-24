// START_MODULE_CONTRACT
//   PURPOSE: Fit compact records, catalog, and help to terminal dimensions.
//   SCOPE: 20/80 records columns, shared compact spacing, full-screen catalog sizing, and footer reservation.
//   DEPENDS: M-APP-SESSIONS, M-APP-PERFORMANCES, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-EXERCISE, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   recordsPaneStyle - compact records/catalog pane padding
//   recordsPaneSpacing - narrow horizontal and zero vertical pane gaps
//   panelSize - calculate one records or catalog allocation
//   resize - resize all tables, page states, dialogs, and help
// END_MODULE_MAP

package views

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_LAYOUT
var recordsPaneStyle = components.BaseStyle

var recordsPaneSpacing = components.PanelSpacing

func (m RecordsModel) panelSize(p panel) (int, int) {
	if p == catalogPanel {
		height := m.height
		if status := m.catalogStatusView(); status != "" {
			height -= lipgloss.Height(status)
		}
		return m.width, max(
			0,
			height,
		)
	}

	usable := max(
		0,
		m.height-lipgloss.Height(m.helpView()),
	)
	fifth := m.width / 5
	if p == sessionsPanel {
		return fifth, usable
	}

	right := m.width - fifth
	if p == performancesPanel {
		return right, usable / 2
	}

	return right, usable - usable/2
}

func (m RecordsModel) resize() RecordsModel {
	tables := []*table.Model{&m.sessionTable, &m.performanceTable, &m.setTable, &m.catalogTable}
	for i, t := range tables {
		outerWidth, outerHeight := m.panelSize(panel(i))
		width := max(
			0,
			outerWidth-recordsPaneStyle.GetHorizontalFrameSize()-recordsPaneSpacing.GetHorizontalFrameSize(),
		)
		height := max(
			1,
			outerHeight-recordsPaneStyle.GetVerticalFrameSize()-recordsPaneSpacing.GetVerticalFrameSize(),
		)
		t.SetColumns(components.EvenColumns(
			t.Columns(),
			width,
		))
		t.SetWidth(width)
		t.SetHeight(height)
	}
	m.syncPagers()
	if m.dialog != nil {
		m.dialog.Resize(
			m.width,
			m.height,
		)
	}

	if m.showingHelp {
		m.help.resize(
			m.width,
			m.height,
		)
	}

	return m
}

// END_BLOCK_RECORDS_LAYOUT
