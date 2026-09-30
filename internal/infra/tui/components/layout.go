package components

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

func EvenColumns(source []table.Column, width int) []table.Column {
	columns := append([]table.Column(nil), source...)
	if len(columns) > 0 {
		// Column widths exclude the padding applied to each cell.
		cellPadding := table.DefaultStyles().Cell.GetHorizontalFrameSize()
		for i := range columns {
			columnWidth := width / len(columns)
			if i < width%len(columns) {
				columnWidth++
			}

			columns[i].Width = max(0, columnWidth-cellPadding)
		}
	}

	return columns
}

func Position(cursor, count int) string {
	if count == 0 {
		return "0 of 0"
	}

	return fmt.Sprintf("%d of %d", cursor+1, count)
}

func RenderPane(style lipgloss.Style, title, content, counter string) string {
	rendered := style.Render(content)
	lines := strings.Split(rendered, "\n")
	border := style.GetBorderStyle()
	width := max(0, lipgloss.Width(lines[0])-2)

	label := []rune(" " + title + " ")
	if len(label) > width {
		label = label[:width]
	}

	topStyle := lipgloss.NewStyle().Foreground(style.GetBorderTopForeground())
	lines[0] = topStyle.Render(border.TopLeft) + string(label) +
		topStyle.Render(strings.Repeat(border.Top, width-len(label))+border.TopRight)

	footer := []rune(" " + counter + " ")
	if len(footer) > width {
		footer = footer[:width]
	}

	bottom := border.BottomLeft + strings.Repeat(border.Bottom, width-len(footer)) + string(footer) + border.BottomRight
	lines[len(lines)-1] = lipgloss.NewStyle().Foreground(style.GetBorderBottomForeground()).Render(bottom)

	return strings.Join(lines, "\n")
}
