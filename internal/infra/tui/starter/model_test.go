package starter

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/danielsavinoff/trainings-logger/internal/app/command"
	"github.com/danielsavinoff/trainings-logger/internal/app/query"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence/sqlite"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"

	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/testutils"
)

type childSpy struct {
	calls int
}

func (s *childSpy) Init() tea.Cmd {
	return nil
}

func (s *childSpy) View() tea.View {
	return tea.NewView("")
}

func (s *childSpy) Update(tea.Msg) (tea.Model, tea.Cmd) {
	s.calls++
	return s, nil
}

func TestGlobalQuitShortcuts(t *testing.T) {
	for _, key := range []string{"ctrl+c", "ctrl+z"} {
		child := &childSpy{}
		m := Model{child: child}
		_, cmd := m.Update(testutils.Key(key))
		if cmd == nil {
			t.Fatalf(
				"%s did not quit",
				key,
			)
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatalf(
				"%s must quit, not suspend",
				key,
			)
		}
		if child.calls != 0 {
			t.Fatal("global quit reached child")
		}
	}
	child := &childSpy{}
	m := Model{child: child}
	for _, key := range []string{"q", "?", "j"} {
		_, cmd := m.Update(testutils.Key(key))
		if cmd != nil {
			t.Fatal("ordinary input should be delegated")
		}
	}
	if child.calls != 3 {
		t.Fatal("ordinary input did not reach child")
	}
}

func testApplication(t *testing.T) Model {
	t.Helper()
	ctx := context.Background()
	db, err := sqlite.OpenFile(
		ctx,
		filepath.Join(
			t.TempDir(),
			"test.db",
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	exercises := sqlite.NewExerciseRepository(db)
	sets := sqlite.NewSetRepository(db)
	prefs := sqlite.NewPreferencesRepository(db)
	app := tui.App{
		ListExercises: query.NewListExercisesHandler(
			exercises,
			persistence.NewTransactionRunner[query.ListExercisesResult](db),
		),
		ListSets: query.NewListSetsHandler(
			sets,
			persistence.NewTransactionRunner[query.ListSetsResult](db),
		),
		GetPreferences: query.NewGetPreferencesHandler(
			prefs,
			persistence.NewTransactionRunner[query.GetPreferencesResult](db),
		),
		CreateExercise: command.NewCreateExerciseHandler(
			exercises,
			persistence.NewTransactionRunner[command.CreateExerciseResult](db),
		),
	}
	if _, err := app.CreateExercise.Handle(
		ctx,
		command.CreateExercise{Name: "Squat"},
	); err != nil {
		t.Fatal(err)
	}
	m, err := NewModel(
		ctx,
		app,
	)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func appPress(m Model, keys ...string) Model {
	for _, key := range keys {
		next, _ := m.Update(testutils.Key(key))
		m = next.(Model)
	}
	return m
}

func TestApplicationQuitContexts(t *testing.T) {
	for _, screen := range []struct {
		name string
		keys []string
	}{
		{
			name: "records",
		},
		{
			name: "help",
			keys: []string{"?"},
		},
		{
			name: "dialog",
			keys: []string{"n", "i"},
		},
	} {
		for _, key := range []string{"ctrl+c", "ctrl+z", "q"} {
			t.Run(
				screen.name+"/"+key,
				func(t *testing.T) {
					m := appPress(
						testApplication(t),
						screen.keys...,
					)
					next, cmd := m.Update(testutils.Key(key))
					if screen.name == "help" {
						if cmd != nil {
							t.Fatal("a key in help should return without quitting")
						}
						if next.(Model).child.(interface{ HelpActive() bool }).HelpActive() {
							t.Fatal("key did not return to records")
						}
						return
					}
					if screen.name == "dialog" && key == "q" {
						if cmd != nil {
							if _, ok := cmd().(tea.QuitMsg); ok {
								t.Fatal("q quit during text entry")
							}
						}
						m = appPress(
							next.(Model),
							"?",
						)
						if !strings.Contains(
							ansi.Strip(m.View().Content),
							"q?",
						) {
							t.Fatal("dialog did not retain q and ? as text")
						}
						return
					}
					if cmd == nil {
						t.Fatal("quit shortcut returned no command")
					}
					if _, ok := cmd().(tea.QuitMsg); !ok {
						t.Fatal("shortcut did not quit")
					}
				},
			)
		}
	}
}
