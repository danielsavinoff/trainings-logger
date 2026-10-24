// START_MODULE_CONTRACT
//   PURPOSE: Render compact records hints and action outcomes.
//   SCOPE: Align a single help hint and actual transfer/error feedback with table spacing.
//   DEPENDS: none
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   footerLine - wrap footer text inside the shared horizontal panel spacing
//   helpView - render only the help shortcut plus real status/error feedback
// END_MODULE_MAP

package views

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_FOOTER
func (m RecordsModel) footerLine(text string) string {
	width := max(
		1,
		m.width-recordsPaneSpacing.GetHorizontalFrameSize(),
	)
	return lipgloss.NewStyle().Width(width).MaxWidth(width).Render(
		ansi.Hardwrap(
			ansi.Wrap(
				text,
				width,
				"",
			),
			width,
			false,
		))
}

func (m RecordsModel) helpView() string {
	lines := []string{}
	if m.transferStatus != "" {
		lines = append(
			lines,
			m.footerLine(m.transferStatus),
		)
	}
	if m.err != "" {
		lines = append(
			lines,
			m.footerLine(m.err),
		)
	}
	hint := lipgloss.NewStyle().Foreground(components.BaseStyle.GetBorderTopForeground()).Render(m.footerLine("? - help"))
	lines = append(
		lines,
		hint,
	)
	return recordsPaneSpacing.Render(strings.Join(
		lines,
		"\n",
	))
}

// END_BLOCK_RECORDS_FOOTER
