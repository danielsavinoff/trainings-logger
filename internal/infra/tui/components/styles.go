// START_MODULE_CONTRACT
//   PURPOSE: Share compact frame padding and outer spacing across every TUI view.
//   SCOPE: Shared colors, one-cell horizontal padding, and zero vertical padding.
//   DEPENDS: none
//   LINKS: M-TUI-COMPONENTS, V-M-TUI-COMPONENTS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   BaseStyle - define common bordered pane padding
//   PanelSpacing - define common outer spacing
//   SelectedStyle - style selected table rows
// END_MODULE_MAP

package components

import (
	"charm.land/lipgloss/v2"
)

// START_BLOCK_SHARED_STYLES

var ActivePanelColor = lipgloss.Color("#8BA4D4")

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("245")).
	Padding(
		0,
		1,
	)

var PanelSpacing = lipgloss.NewStyle().Padding(
	0,
	1,
)

func SelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#182235")).Background(ActivePanelColor)
}

// END_BLOCK_SHARED_STYLES
