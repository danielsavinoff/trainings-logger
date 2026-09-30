package components

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/paginator"
	"charm.land/lipgloss/v2"
)

func NewPaginator() paginator.Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(ActivePanelColor).Render("● ")
	p.InactiveDot = lipgloss.NewStyle().Foreground(BaseStyle.GetBorderTopForeground()).Render("○ ")
	p.KeyMap.PrevPage = key.NewBinding(key.WithKeys("left", "h", "pgup", "b"))
	p.KeyMap.NextPage = key.NewBinding(key.WithKeys("right", "l", "pgdown", "f", "space"))

	return p
}

func PaginatorView(p paginator.Model, width int) string {
	// Long histories use numbers so the indicator always fits on one line.
	if p.TotalPages*2 > width {
		p.Type = paginator.Arabic
	}

	return lipgloss.NewStyle().Width(width).MaxWidth(width).Align(lipgloss.Center).Render(p.View())
}
