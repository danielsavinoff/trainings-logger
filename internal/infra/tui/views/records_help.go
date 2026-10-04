// START_MODULE_CONTRACT
//   PURPOSE: Render compact records footer and errors.
//   SCOPE: Wrap help/quit hints and preserve visible status text.
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
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
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
	hints := "?  help    q  quit"
	if lipgloss.Width(hints) > width {
		hints = "?  help\nq  quit"
	}
	status := ""
	if m.err != "" {
		status = center.Render(wrap(m.err)) + "\n"
	}
	return status + center.Render(lipgloss.NewStyle().Foreground(components.BaseStyle.GetBorderTopForeground()).Render(wrap(hints)))
}

// END_BLOCK_RECORDS_FOOTER
