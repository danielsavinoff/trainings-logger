// START_MODULE_CONTRACT
//   PURPOSE: Render records, help, or setup in the alternate screen.
//   SCOPE: Terminal modes, panels, dialog layers, and cursor.
//   DEPENDS: M-TUI-HELP, M-TUI-QUESTIONNAIRE
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS, M-TUI-QUESTIONNAIRE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   View - compose the active screen and terminal modes
// END_MODULE_MAP

package views

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_VIEW
func (m RecordsModel) View() tea.View {
	if m.questionnaire != nil {
		return m.questionnaireView()
	}
	if m.showingHelp {
		v := tea.NewView(m.help.view())
		v.AltScreen = true
		v.MouseMode = tea.MouseModeCellMotion
		return v
	}
	styles := []lipgloss.Style{components.BaseStyle, components.BaseStyle}
	styles[m.activePanel] = styles[m.activePanel].BorderForeground(components.ActivePanelColor)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			components.PanelSpacing.Render(components.RenderPane(
				styles[exercisesPanel],
				"Exercises <1>",
				m.exercisesView(),
				components.Position(
					m.table.Cursor(),
					len(m.table.Rows()),
				),
			)),
			components.PanelSpacing.Render(components.RenderPane(
				styles[setsPanel],
				"Sets <2>",
				m.setsView(),
				components.Position(
					m.setCursor,
					len(m.setLines),
				),
			)),
		),
		m.helpView(),
	)
	if m.width > 0 && m.height > 0 {
		content = lipgloss.NewStyle().MaxWidth(m.width).MaxHeight(m.height).Render(content)
	}

	layers := []*lipgloss.Layer{lipgloss.NewLayer(content).ID("panels")}
	var cursor *tea.Cursor
	if m.dialog != nil {
		var dialogLayers []*lipgloss.Layer
		dialogLayers, cursor = m.dialog.Layers()
		layers = append(
			layers,
			dialogLayers...,
		)
	}

	canvas := lipgloss.NewCanvas(
		max(
			0,
			m.width,
		),
		max(
			0,
			m.height,
		),
	)
	content = canvas.Compose(lipgloss.NewCompositor(layers...)).Render()

	v := tea.NewView(content)
	v.Cursor = cursor
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	return v
}

// END_BLOCK_RECORDS_VIEW
