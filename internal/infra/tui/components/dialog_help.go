package components

import (
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

func (d *Dialog) helpView(width int) string {
	bindings := []key.Binding{
		d.keys.Up, d.keys.Down, d.keys.Cancel, d.keys.Insert,
	}

	help := help.New()
	help.SetWidth(0)
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)

	var lines []string
	for _, row := range [][]key.Binding{bindings[:2], bindings[2:]} {
		content := help.ShortHelpView(row)
		if lipgloss.Width(content) <= width {
			lines = append(lines, center.Render(content))
			continue
		}

		for _, binding := range row {
			lines = append(lines, center.Render(help.ShortHelpView([]key.Binding{binding})))
		}
	}

	return strings.Join(lines, "\n")
}
