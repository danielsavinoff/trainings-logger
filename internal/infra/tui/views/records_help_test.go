package views

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/danielsavinoff/trainings-logger/internal/app/command"
	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/app/query"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/testutils"
)

func TestRecordsHelpStatePreservation(t *testing.T) {
	for _, keys := range [][]string{{"j"}, {"g", "s", "l", "j"}} {
		for _, back := range []string{"esc", "?", "q", "ctrl+c", "ctrl+z", "j", "n"} {
			m := press(
				sampleModel(),
				keys...,
			)
			m.err = "An existing error"
			before := m.View().Content
			panel, exercise, set := m.activePanel, m.table.Cursor(), m.setCursor
			page, offset := m.setPager.Page, m.timeline.YOffset()
			m = press(
				m,
				"?",
			)
			if !m.showingHelp {
				t.Fatal("help did not open")
			}
			m = press(
				m,
				back,
			)
			if m.showingHelp || m.View().Content != before || m.activePanel != panel || m.table.Cursor() != exercise || m.setCursor != set || m.setPager.Page != page || m.timeline.YOffset() != offset {
				t.Fatalf(
					"help changed browsing state for %v via %s",
					keys,
					back,
				)
			}
			m = press(
				m,
				"?",
			)
			updated, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
			m = press(
				updated.(RecordsModel),
				back,
			)
			if m.activePanel != panel || m.table.Cursor() != exercise || m.setCursor != set || m.err != "An existing error" {
				t.Fatal("resizing help lost selection, focus or error")
			}
			if len(m.setLines) > 0 && m.setPager.Page != m.setLines[m.setCursor]/m.setPager.PerPage {
				t.Fatal("resized pagination is inconsistent")
			}
		}
	}
}

type transactionSpy[T any] struct {
	calls *int
}

func (s transactionSpy[T]) Run(context.Context, func(context.Context) (T, error)) (T, error) {
	*s.calls++
	var result T
	return result, errors.New("unexpected persistence access")
}

func TestRecordsHelpInputIsolation(t *testing.T) {
	for _, panel := range []panel{exercisesPanel, setsPanel} {
		m := sampleModel()
		m.focusPanel(panel)
		calls := 0
		writes := transactionSpy[ports.TransactionNoResult]{&calls}
		m.store = &recordStore{App: tui.App{
			CreateExercise: command.NewCreateExerciseHandler(
				nil,
				transactionSpy[command.CreateExerciseResult]{&calls},
			),
			CreateSet: command.NewCreateSetHandler(
				nil,
				nil,
				nil,
				transactionSpy[command.CreateSetResult]{&calls},
			),
			UpdateExercise: command.NewUpdateExerciseHandler(
				nil,
				writes,
			),
			UpdateSet: command.NewUpdateSetHandler(
				nil,
				writes,
			),
			DeleteExercise: command.NewDeleteExerciseHandler(
				nil,
				writes,
			),
			DeleteSet: command.NewDeleteSetHandler(
				nil,
				writes,
			),
			ListSets: query.NewListSetsHandler(
				nil,
				transactionSpy[query.ListSetsResult]{&calls},
			),
		}}
		before := m.View().Content
		exercises := append(
			[]exerciseRecord(nil),
			m.exercises...,
		)
		m = press(
			m,
			"g",
			"?",
		)
		if !m.showingHelp || m.goPending {
			t.Fatal("help did not clear pending focus prefix")
		}
		for _, key := range []string{"n", "e", "d", "enter", "g", "s", "j", "l", "q", "esc", "?", "ctrl+c", "ctrl+z"} {
			probe := m
			updated, cmd := probe.Update(testutils.Key(key))
			probe = updated.(RecordsModel)
			if probe.showingHelp {
				t.Fatalf(
					"%s did not close help",
					key,
				)
			}
			if key == "ctrl+c" || key == "ctrl+z" || key == "q" || key == "esc" || key == "?" {
				if cmd != nil {
					t.Fatalf(
						"%s performed an action while dismissing help",
						key,
					)
				}
			} else if key == "n" || key == "e" || key == "d" || key == "enter" {
				if cmd != nil || probe.dialog != nil || !reflect.DeepEqual(
					exercises,
					probe.exercises,
				) {
					t.Fatalf(
						"%s triggered its records action",
						key,
					)
				}
			}
		}

		updated, _ := m.Update(tea.MouseWheelMsg{X: 60, Y: 6, Button: tea.MouseWheelDown})
		m = press(
			updated.(RecordsModel),
			"esc",
		)
		if calls != 0 || m.dialog != nil || !reflect.DeepEqual(
			exercises,
			m.exercises,
		) || m.View().Content != before {
			t.Fatal("help input changed records or accessed persistence")
		}
	}
}

func TestRecordsHelpQuitAndDialogInput(t *testing.T) {
	for _, help := range []bool{false} {
		m := sampleModel()
		if help {
			m = press(
				m,
				"?",
			)
		}
		_, cmd := m.Update(testutils.Key("q"))
		if cmd == nil {
			t.Fatal("q did not quit")
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatal("q did not return quit command")
		}
	}
	m := press(
		sampleModel(),
		"n",
		"i",
		"?",
		"q",
	)
	if m.showingHelp || m.dialog == nil || m.dialog.Value(0) != "?q" {
		t.Fatal("dialog text intercepted")
	}
	m = press(
		m,
		"esc",
		"esc",
		"?",
	)
	if !m.showingHelp || m.View().Cursor != nil || !m.View().AltScreen {
		t.Fatal("help view modes or cursor incorrect")
	}
}

func TestRecordsCompactFooter(t *testing.T) {
	for _, width := range []int{100, 40, 12, 8, 1} {
		m := sampleModel()
		m.width = width
		styledFooter := m.helpView()
		if strings.Contains(
			styledFooter,
			"38;2;139;164;212m",
		) || strings.Contains(
			ansi.Strip(styledFooter),
			"╭",
		) {
			t.Fatal("footer must be borderless and use inactive border color")
		}
		footer := ansi.Strip(styledFooter)
		if lipgloss.Width(footer) > width {
			t.Fatalf(
				"footer exceeds width %d",
				width,
			)
		}
		if width >= 8 && (!strings.Contains(
			footer,
			"?  help",
		) || !strings.Contains(
			footer,
			"q  quit",
		)) {
			t.Fatal("footer hides hints")
		}
		if width >= 40 && lipgloss.Height(footer) != 1 {
			t.Fatal("footer should fit on one line")
		}
		m.err = "Could not save this record"
		m = m.resize()
		if width >= 40 && !strings.Contains(
			ansi.Strip(m.View().Content),
			m.err,
		) {
			t.Fatal("footer lost error")
		}
	}
}
