// START_MODULE_CONTRACT
//   PURPOSE: Load and format sets belonging to the selected performance.
//   SCOPE: Domain set reads and display-only weight conversion without date separators.
//   DEPENDS: M-APP-SESSIONS, M-DOMAIN-EXERCISE, M-DOMAIN-SESSION, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   selectedSet - locate the selected set
//   loadSets - query the selected performance and retain set selection
//   refreshSetRows - format reps, weight, and rest
//   refreshSelectedSets - refresh children after session navigation
// END_MODULE_MAP

package views

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/table"

	"github.com/danielsavinoff/trainings-logger/internal/app/query"
)

// START_BLOCK_RECORDS_SETS
func (m RecordsModel) selectedSet() int {
	i := m.setTable.Cursor()
	if i < 0 || i >= len(m.sets) {
		return -1
	}

	return i
}

func (m *RecordsModel) loadSets(ctx context.Context, preferred string) error {
	performanceID := ""
	if i := m.selectedPerformance(); i >= 0 {
		performanceID = m.performances[i].ID
	}

	fallback := 0
	if performanceID == m.setsPerformanceID {
		fallback = m.setTable.Cursor()
	}

	m.setsPerformanceID = performanceID
	m.sets = nil
	m.setTable.SetRows(nil)
	m.syncPagers()
	if performanceID == "" || m.store == nil {
		return nil
	}
	result, err := m.store.ListSets.Handle(
		ctx,
		query.ListSets{PerformanceID: performanceID},
	)
	if err != nil {
		return err
	}

	m.sets = result.Result
	if err := m.refreshSetRows(); err != nil {
		return err
	}

	ids := make(
		[]string,
		0,
		len(m.sets),
	)
	for _, item := range m.sets {
		ids = append(
			ids,
			item.ID,
		)
	}
	restoreCursor(
		&m.setTable,
		ids,
		preferred,
		fallback,
	)
	m.syncPagers()
	return nil
}

func (m *RecordsModel) refreshSetRows() error {
	rows := make(
		[]table.Row,
		0,
		len(m.sets),
	)
	for _, item := range m.sets {
		weight, err := item.ExerciseWeight.Weight.To(m.weightUnit)
		if err != nil {
			return err
		}

		rows = append(
			rows,
			table.Row{fmt.Sprint(item.Repetitions), fmt.Sprintf(
				"%.5g %s",
				weight.Value,
				strings.ToUpper(string(m.weightUnit)),
			), fmt.Sprint(item.RestInSeconds)},
		)
	}

	cursor := m.setTable.Cursor()
	m.setTable.SetRows(rows)
	m.setTable.SetCursor(cursor)
	return nil
}

func (m *RecordsModel) refreshSelectedSets() {
	if err := m.loadSets(
		context.Background(),
		"",
	); err != nil {
		m.err = err.Error()
	} else {
		m.err = ""
	}
}

// END_BLOCK_RECORDS_SETS
