package views

import (
	"fmt"

	"charm.land/bubbles/v2/table"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

func newExerciseTable(exercises []exerciseRecord) table.Model {
	columns := []table.Column{{Title: "Name"}, {Title: "Body Weight %"}}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(exerciseRows(exercises)),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(42),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(components.ActivePanelColor).
		Bold(false)
	s.Selected = components.SelectedStyle()
	t.SetStyles(s)

	return t
}

func exerciseRows(exercises []exerciseRecord) []table.Row {
	rows := make([]table.Row, 0, len(exercises))
	for _, exercise := range exercises {
		percentage := fmt.Sprintf("%.5g%%", float64(exercise.bodyWeightShare)*100)
		rows = append(rows, table.Row{exercise.name, percentage})
	}

	return rows
}

func (m RecordsModel) selectedExercise() int {
	i := m.table.Cursor()
	if i < 0 || i >= len(m.exercises) {
		return -1
	}

	return i
}

func (m *RecordsModel) syncExercisePager() {
	m.exercisePager.PerPage = max(1, m.table.Height())
	m.exercisePager.SetTotalPages(max(1, len(m.table.Rows())))
	m.exercisePager.Page = max(0, m.table.Cursor()) / m.exercisePager.PerPage
}

func (m RecordsModel) exercisesView() string {
	start, end := m.exercisePager.GetSliceBounds(len(m.table.Rows()))
	t := m.table
	t.SetRows(m.table.Rows()[start:end])
	t.GotoTop()
	t.SetCursor(max(0, m.table.Cursor()-start))

	return t.View()
}
