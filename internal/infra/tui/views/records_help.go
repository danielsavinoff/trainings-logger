// START_MODULE_CONTRACT
//   PURPOSE: Render compact records hints and action outcomes.
//   SCOPE: Wrap help/load/save/quit hints, transfer progress, and error text.
//   DEPENDS: none
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   helpView - render records command hints and errors
// END_MODULE_MAP

package views

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_FOOTER
func (m RecordsModel) helpView() string {
	width := max(
		1,
		m.width,
	)
	left := lipgloss.NewStyle().Width(width).Align(lipgloss.Left)
	wrap := func(text string) string {
		return ansi.Hardwrap(
			ansi.Wrap(
				text,
				width,
				"",
			),
			width,
			false,
		)
	}
	hints := "? - help  l - load  s - save  q - quit"
	if lipgloss.Width(hints) > width {
		hints = "? - help\nl - load\ns - save\nq - quit"
	}
	status := ""
	if m.transferStatus != "" {
		status = left.Render(wrap(m.transferStatus)) + "\n"
	}
	if m.err != "" {
		status += left.Render(wrap(m.err)) + "\n"
	}
	return status + left.Render(lipgloss.NewStyle().Foreground(components.BaseStyle.GetBorderTopForeground()).Render(wrap(hints)))
}

// END_BLOCK_RECORDS_FOOTER
