package components

import (
	"charm.land/lipgloss/v2"
)

var ActivePanelColor = lipgloss.Color("#8BA4D4")

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("245")).
	Padding(1, 2)

var PanelSpacing = lipgloss.NewStyle().Padding(1, 2)

func SelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#182235")).Background(ActivePanelColor)
}
