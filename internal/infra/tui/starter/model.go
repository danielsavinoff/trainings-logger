// Package starter composes the TUI and handles application-wide input.
package starter

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/views"
)

var _ tea.Model = Model{}

type Model struct {
	child tea.Model
}

// NewModel loads the records view and wraps it in the application model.
func NewModel(ctx context.Context, app tui.App) (Model, error) {
	child, err := views.NewRecordsModel(ctx, app)
	if err != nil {
		return Model{}, err
	}

	return Model{child: child}, nil
}

func (m Model) Init() tea.Cmd { return m.child.Init() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if press, ok := msg.(tea.KeyPressMsg); ok && press.String() == "ctrl+c" {
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View { return m.child.View() }
