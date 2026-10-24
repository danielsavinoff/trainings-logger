// START_MODULE_CONTRACT
//   PURPOSE: Render independent catalog rows and performance rows.
//   SCOPE: Shared table styling and consistent catalog/workout exercise percentage columns.
//   DEPENDS: M-APP-SESSIONS, M-APP-PERFORMANCES, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-EXERCISE, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   newRecordsTable - construct a styled records table
//   exerciseRow - format one catalog definition with body-weight percentage
//   catalogRows - format reusable exercise names and percentages
//   selectedCatalogExercise - locate the catalog selection
//   selectedPerformance - locate the selected workout performance
//   performanceRows - resolve performance names and percentages from catalog definitions
// END_MODULE_MAP

package views

import (
	"fmt"

	"charm.land/bubbles/v2/table"

	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_EXERCISES
func newRecordsTable(columns []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
		table.WithWidth(42),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.Foreground(components.ActivePanelColor).Bold(false)
	s.Selected = components.SelectedStyle()
	t.SetStyles(s)
	return t
}

func exerciseRow(item exercise.Exercise) table.Row {
	return table.Row{item.Name, fmt.Sprintf(
		"%.5g%%",
		float64(item.BodyWeightShare)*100,
	)}
}

func catalogRows(catalog []exercise.Exercise) []table.Row {
	rows := make(
		[]table.Row,
		0,
		len(catalog),
	)
	for _, item := range catalog {
		rows = append(
			rows,
			exerciseRow(item),
		)
	}

	return rows
}

func (m RecordsModel) selectedCatalogExercise() int {
	i := m.catalogTable.Cursor()
	if i < 0 || i >= len(m.catalog) {
		return -1
	}

	return i
}

func (m RecordsModel) selectedPerformance() int {
	i := m.performanceTable.Cursor()
	if i < 0 || i >= len(m.performances) {
		return -1
	}

	return i
}

func (m RecordsModel) performanceRows() []table.Row {
	byID := make(
		map[string]table.Row,
		len(m.catalog),
	)
	for _, item := range m.catalog {
		byID[item.ID] = exerciseRow(item)
	}
	rows := make(
		[]table.Row,
		0,
		len(m.performances),
	)
	for _, item := range m.performances {
		row, ok := byID[item.ExerciseID]
		if !ok {
			row = table.Row{"", ""}
		}
		rows = append(
			rows,
			row,
		)
	}
	return rows
}

// END_BLOCK_RECORDS_EXERCISES
