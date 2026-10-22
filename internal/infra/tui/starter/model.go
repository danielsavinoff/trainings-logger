// START_MODULE_CONTRACT
//   PURPOSE: Compose the TUI and handle global quit shortcuts.
//   SCOPE: Application startup and Ctrl+C/Ctrl+Z dispatch.
//   DEPENDS: M-TUI-RECORDS
//   LINKS: M-TUI-STARTER, V-M-TUI-STARTER
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   Model - wrap the current TUI child
//   NewModel - load initial records
//   Init - initialize child
//   Update - quit globally or delegate input
//   View - render child
// END_MODULE_MAP

package starter

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/views"
)

// START_BLOCK_GLOBAL_INPUT
var _ tea.Model = Model{}

type Model struct {
	child tea.Model
}

func NewModel(ctx context.Context, app tui.App) (Model, error) {
	child, err := views.NewRecordsModel(ctx, app)
	if err != nil {
		return Model{}, err
	}

	return Model{child: child}, nil
}

func (m Model) Init() tea.Cmd { return m.child.Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if press, ok := msg.(tea.KeyPressMsg); ok && press.String() != "" {
		if active, ok := m.child.(interface{ HelpActive() bool }); !ok || !active.HelpActive() {
			switch press.String() {
			case "ctrl+c", "ctrl+z":
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View { return m.child.View() }

// END_BLOCK_GLOBAL_INPUT
