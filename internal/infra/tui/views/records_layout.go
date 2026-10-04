// START_MODULE_CONTRACT
//   PURPOSE: Fit records panels and help to terminal dimensions.
//   SCOPE: Reserve footer and error space and propagate resizing.
//   DEPENDS: M-TUI-HELP
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   resize - update panel geometry, pagination, and active help
// END_MODULE_MAP

package views

import (
	"charm.land/lipgloss/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_LAYOUT
func (m RecordsModel) resize() RecordsModel {
	helpHeight := lipgloss.Height(m.helpView())

	frameWidth := components.BaseStyle.GetHorizontalFrameSize() + components.PanelSpacing.GetHorizontalFrameSize()
	width := max(
		0,
		m.width/2-frameWidth,
	)
	rightWidth := max(
		0,
		m.width-m.width/2-frameWidth,
	)
	panelHeight := max(
		0,
		m.height-helpHeight-components.PanelSpacing.GetVerticalFrameSize(),
	)
	frameHeight := components.BaseStyle.GetVerticalFrameSize()

	m.table.SetColumns(components.EvenColumns(
		m.table.Columns(),
		width,
	))
	m.table.SetWidth(width)

	// Reserve at least one line for the column names.
	m.table.SetHeight(max(
		1,
		panelHeight-frameHeight-1,
	))

	m.timeline.SetWidth(rightWidth)

	// Keep the Sets column headings fixed above the scrolling history.
	m.timeline.SetHeight(max(
		0,
		panelHeight-frameHeight-2,
	))

	m.syncExercisePager()
	m.refreshTimeline()

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
