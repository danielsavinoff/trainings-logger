package views

import (
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

func (m RecordsModel) helpView() string {
	bindings := []key.Binding{
		key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
		key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
		key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "move left")),
		key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "move right")),
	}

	actions := []key.Binding{
		m.keys.New, m.keys.Edit, m.keys.Delete, m.keys.Quit,
	}

	help := m.table.Help
	help.SetWidth(m.width)
	center := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center)

	status := ""
	if m.err != "" {
		status = center.Render(m.err) + "\n"
	}

	return status + lipgloss.JoinVertical(lipgloss.Left,
		center.Render(help.ShortHelpView(bindings)),
		center.Render(help.ShortHelpView(actions)),
	)
}
