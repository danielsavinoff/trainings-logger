// START_MODULE_CONTRACT
//   PURPOSE: Load exercises belonging to the selected workout.
//   SCOPE: Workout selection, retained performance IDs, and dependent set loading.
//   DEPENDS: M-APP-SESSIONS, M-APP-PERFORMANCES, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-EXERCISE, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   selectedSession - locate the selected workout
//   restoreCursor - retain an ID or clamp a fallback cursor
//   loadPerformances - query workout membership and selected sets
//   refreshSelectedPerformances - refresh children after workout navigation
// END_MODULE_MAP

package views

import (
	"context"

	"charm.land/bubbles/v2/table"

	"github.com/danielsavinoff/trainings-logger/internal/app/query"
)

// START_BLOCK_RECORDS_SESSIONS
func (m RecordsModel) selectedSession() int {
	i := m.sessionTable.Cursor()
	if i < 0 || i >= len(m.sessions) {
		return -1
	}

	return i
}

func restoreCursor(t *table.Model, ids []string, preferred string, fallback int) {
	cursor := min(
		max(
			0,
			fallback,
		),
		max(
			0,
			len(ids)-1,
		),
	)
	for i, id := range ids {
		if id == preferred && preferred != "" {
			cursor = i
			break
		}
	}
	t.SetCursor(cursor)
}

func (m *RecordsModel) loadPerformances(ctx context.Context, preferred string, setID string) error {
	sessionID := ""
	if i := m.selectedSession(); i >= 0 {
		sessionID = m.sessions[i].ID
	}

	fallback := 0
	if sessionID == m.performancesSessionID {
		fallback = m.performanceTable.Cursor()
	}

	m.performancesSessionID = sessionID
	m.performances = nil
	m.sets = nil
	m.performanceTable.SetRows(nil)
	m.syncPagers()
	if sessionID == "" || m.store == nil {
		m.setTable.SetRows(nil)
		m.syncPagers()
		return nil
	}
	result, err := m.store.ListPerformances.Handle(
		ctx,
		query.ListPerformances{SessionID: sessionID},
	)
	if err != nil {
		m.setTable.SetRows(nil)
		m.syncPagers()
		return err
	}

	m.performances = result.Result
	m.performanceTable.SetRows(m.performanceRows())
	ids := make(
		[]string,
		0,
		len(m.performances),
	)
	for _, item := range m.performances {
		ids = append(
			ids,
			item.ID,
		)
	}
	restoreCursor(
		&m.performanceTable,
		ids,
		preferred,
		fallback,
	)
	m.syncPagers()
	return m.loadSets(
		ctx,
		setID,
	)
}

func (m *RecordsModel) refreshSelectedPerformances() {
	if err := m.loadPerformances(
		context.Background(),
		"",
		"",
	); err != nil {
		m.err = err.Error()
	} else {
		m.err = ""
	}
}

// END_BLOCK_RECORDS_SESSIONS
