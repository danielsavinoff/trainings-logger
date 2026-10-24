// START_MODULE_CONTRACT
//   PURPOSE: Render the workout workspace, help, or setup.
//   SCOPE: Three compact records panes or a full-screen catalog, footer, modal layers, and terminal cursor.
//   DEPENDS: M-APP-SESSIONS, M-APP-PERFORMANCES, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-EXERCISE, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS, M-TUI-QUESTIONNAIRE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   recordsPane - render and bound a single panel
//   View - compose the active alternate-screen view
// END_MODULE_MAP

package views

import (
	"charm.land/bubbles/v2/paginator"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_VIEW
func (m RecordsModel) recordsPane(p panel, title string, t table.Model, pager paginator.Model) string {
	width, height := m.panelSize(p)
	if width <= 0 || height <= 0 {
		return ""
	}

	style := recordsPaneStyle
	if p == m.activePanel {
		style = style.BorderForeground(components.ActivePanelColor)
	}

	pane := recordsPaneSpacing.Render(components.RenderPane(
		style,
		title,
		pageTableView(
			t,
			pager,
		),
		components.Position(
			t.Cursor(),
			len(t.Rows()),
		),
	))
	return lipgloss.NewStyle().MaxWidth(width).MaxHeight(height).Render(pane)
}

func (m RecordsModel) View() tea.View {
	if m.width <= 0 || m.height <= 0 {
		v := tea.NewView("")
		v.AltScreen = true
		return v
	}

	if m.questionnaire != nil {
		return m.questionnaireView()
	}

	if m.showingHelp {
		v := tea.NewView(m.help.view())
		v.AltScreen = true
		v.MouseMode = tea.MouseModeCellMotion
		return v
	}

	var content string
	if m.showingCatalog {
		content = m.catalogView()
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.recordsPane(
					sessionsPanel,
					"Sessions <1>",
					m.sessionTable,
					m.sessionPager,
				),
				lipgloss.JoinVertical(
					lipgloss.Left,
					m.recordsPane(
						performancesPanel,
						"Performances <2>",
						m.performanceTable,
						m.performancePager,
					),
					m.recordsPane(
						setsPanel,
						"Sets <3>",
						m.setTable,
						m.setPager,
					),
				),
			),
			m.helpView(),
		)
	}

	content = lipgloss.NewStyle().MaxWidth(m.width).MaxHeight(m.height).Render(content)
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
		m.width,
		m.height,
	)
	content = canvas.Compose(lipgloss.NewCompositor(
		layers...,
	)).Render()
	v := tea.NewView(content)
	v.Cursor = cursor
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// END_BLOCK_RECORDS_VIEW
