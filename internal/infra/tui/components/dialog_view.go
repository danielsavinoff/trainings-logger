// START_MODULE_CONTRACT
//   PURPOSE: Render compact forms using the same inner padding as records.
//   SCOPE: Consecutive fields, local errors, help and modal borders.
//   DEPENDS: none
//   LINKS: M-TUI-COMPONENTS, V-M-TUI-COMPONENTS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   View - render a bounded compact form
//   renderModalFrame - place modal title and save action in borders
// END_MODULE_MAP

package components

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// START_BLOCK_DIALOG_VIEW

func (d *Dialog) View() string {
	width, compact := d.dialogSize()
	innerWidth := width - BaseStyle.GetHorizontalFrameSize()
	if innerWidth <= 0 || d.height <= BaseStyle.GetVerticalFrameSize() {
		return ""
	}

	title := d.title
	if compact {
		title += fmt.Sprintf(
			" · %d/%d",
			d.active+1,
			len(d.fields),
		)
	}

	var lines []string
	for i, field := range d.fields {
		if compact && i != d.active {
			continue
		}

		label := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
		if i == d.active {
			label = label.Foreground(ActivePanelColor)
		}

		if compact {
			input, _ := fieldView(
				field,
				innerWidth,
			)
			lines = append(
				lines,
				label.Render(field.label),
				input,
			)
		} else {
			input, _ := fieldView(
				field,
				innerWidth-17,
			)
			lines = append(
				lines,
				label.Width(17).Render(field.label)+input,
			)
		}
	}

	if d.err != "" {
		lines = append(
			lines,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#F0A3A3")).Width(innerWidth).Render(d.err),
		)
	}

	lines = append(
		lines,
		d.helpView(innerWidth),
	)

	style := BaseStyle.Border(lipgloss.RoundedBorder()).BorderForeground(ActivePanelColor).
		Foreground(lipgloss.Color("#E6EAF2")).
		Width(width).MaxHeight(d.height)

	return renderModalFrame(
		style,
		title,
		strings.Join(
			lines,
			"\n",
		),
		" Save <enter> ",
	)
}

func renderModalFrame(style lipgloss.Style, title, body, footer string) string {
	rendered := strings.Split(
		style.Render(body),
		"\n",
	)
	border := style.GetBorderStyle()
	borderWidth := max(
		0,
		lipgloss.Width(rendered[0])-2,
	)
	label := []rune(" " + title + " ")
	if len(label) > borderWidth {
		label = label[:borderWidth]
	}

	remaining := borderWidth - lipgloss.Width(string(label))
	borderStyle := lipgloss.NewStyle().Foreground(style.GetBorderTopForeground())
	rendered[0] = borderStyle.Render(border.TopLeft+strings.Repeat(
		border.Top,
		remaining/2,
	)) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Render(string(label)) +
		borderStyle.Render(strings.Repeat(
			border.Top,
			remaining-remaining/2,
		)+border.TopRight)

	footer = footer[:min(
		len(footer),
		borderWidth,
	)]
	remaining = borderWidth - len(footer)
	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	rendered[len(rendered)-1] = borderStyle.Render(border.BottomLeft+strings.Repeat(
		border.Bottom,
		remaining/2,
	)) +
		white.Render(footer) + borderStyle.Render(strings.Repeat(
		border.Bottom,
		remaining-remaining/2,
	)+border.BottomRight)

	return strings.Join(
		rendered,
		"\n",
	)
}

// END_BLOCK_DIALOG_VIEW
