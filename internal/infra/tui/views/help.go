// START_MODULE_CONTRACT
//   PURPOSE: Render and scroll records keyboard help.
//   SCOPE: Help content, responsive layout, and isolated viewport input.
//   DEPENDS: none
//   LINKS: M-TUI-HELP, V-M-TUI-HELP
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   helpModel - retain help viewport and terminal dimensions
//   newHelpModel - initialize responsive help
//   resize - arrange shortcut groups for terminal dimensions
//   update - scroll the help viewport
//   view - render shortcut groups without a page heading
// END_MODULE_MAP

package views

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_HELP_VIEW
type helpEntry struct {
	key         string
	explanation string
}

type helpGroup struct {
	title   string
	entries []helpEntry
}

var recordsHelp = []helpGroup{
	{
		title: "Navigation",
		entries: []helpEntry{
			{key: "↑/k", explanation: "move up"},
			{key: "↓/j", explanation: "move down"},
			{key: "ge", explanation: "focus exercises"},
			{key: "gs", explanation: "focus sets"},
			{key: "enter", explanation: "focus sets from exercises"},
		},
	},

	{
		title: "Pages",
		entries: []helpEntry{
			{key: "←/h/pgup/b", explanation: "previous page"},
			{key: "→/l/pgdown/f/space", explanation: "next page"},
			{key: "u/ctrl+u", explanation: "half page up"},
			{key: "ctrl+d", explanation: "half page down"},
			{key: "home/gg", explanation: "first row"},
			{key: "end/G", explanation: "last row"},
		},
	},

	{
		title: "Records (active panel)",
		entries: []helpEntry{
			{key: "n", explanation: "new record"},
			{key: "e", explanation: "edit record"},
			{key: "d", explanation: "delete record"},
		},
	},

	{
		title: "Help",
		entries: []helpEntry{
			{key: "any key", explanation: "return to records"},
		},
	},
}

type helpModel struct {
	viewport      viewport.Model
	width, height int
}

func newHelpModel(width, height int) helpModel {
	m := helpModel{viewport: viewport.New()}
	m.resize(
		width,
		height,
	)
	return m
}

func (m *helpModel) resize(width, height int) {
	m.width, m.height = max(
		0,
		width,
	), max(
		0,
		height,
	)
	m.viewport.SetWidth(m.width)
	m.viewport.SetHeight(m.height)
	m.viewport.SetContent(m.content())
}

func (m helpModel) content() string {
	gap := 4
	colWidth := (m.width - gap) / 2
	if colWidth >= 40 {
		left := strings.Join(
			[]string{m.renderGroup(
				recordsHelp[0],
				colWidth,
			), m.renderGroup(
				recordsHelp[2],
				colWidth,
			)},
			"\n\n",
		)
		right := strings.Join(
			[]string{m.renderGroup(
				recordsHelp[1],
				colWidth,
			), m.renderGroup(
				recordsHelp[3],
				colWidth,
			)},
			"\n\n",
		)
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			left,
			strings.Repeat(
				" ",
				gap,
			),
			right,
		)
	}
	var groups []string
	for _, group := range recordsHelp {
		groups = append(
			groups,
			m.renderGroup(
				group,
				m.width,
			),
		)
	}
	return strings.Join(
		groups,
		"\n\n",
	)
}

func (m helpModel) renderGroup(group helpGroup, width int) string {
	width = max(
		1,
		width,
	)
	title := lipgloss.NewStyle().Bold(true).Render(ansi.Truncate(
		group.title,
		width,
		"",
	))
	maxKey := 0
	for _, entry := range group.entries {
		maxKey = max(
			maxKey,
			ansi.StringWidth(entry.key),
		)
	}
	tab := maxKey + 4
	var lines = []string{title}
	keyStyle := lipgloss.NewStyle().Foreground(components.ActivePanelColor)
	for _, entry := range group.entries {
		descWidth := width - tab
		key := entry.key
		if descWidth < 1 {
			lines = append(
				lines,
				ansi.Hardwrap(
					key,
					width,
					false,
				),
			)
			descWidth = width
			for _, line := range strings.Split(
				ansi.Hardwrap(
					entry.explanation,
					descWidth,
					false,
				),
				"\n",
			) {
				lines = append(
					lines,
					line,
				)
			}
			continue
		}
		desc := strings.Split(
			ansi.Hardwrap(
				entry.explanation,
				descWidth,
				false,
			),
			"\n",
		)
		padding := strings.Repeat(
			" ",
			max(
				0,
				tab-ansi.StringWidth(key),
			),
		)
		lines = append(
			lines,
			keyStyle.Render(key)+padding+desc[0],
		)
		for _, line := range desc[1:] {
			lines = append(
				lines,
				strings.Repeat(
					" ",
					tab,
				)+line,
			)
		}
	}
	return strings.Join(
		lines,
		"\n",
	)
}

func (m helpModel) update(msg tea.Msg) helpModel {
	m.viewport, _ = m.viewport.Update(msg)
	return m
}

func (m helpModel) view() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	return m.viewport.View()
}

// END_BLOCK_HELP_VIEW
