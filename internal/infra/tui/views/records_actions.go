// START_MODULE_CONTRACT
//   PURPOSE: Reload workout records and apply catalog or active-panel actions.
//   SCOPE: Stable four-part selection, idempotent catalog addition, and scoped deletion.
//   DEPENDS: M-APP-SESSIONS, M-APP-PERFORMANCES, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-EXERCISE, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   selectionIDs - identify selected workout, catalog, performance, and set
//   reloadRecords - refresh independent root lists and selected children
//   addCatalogExercise - add/select a catalog exercise and return to records
//   createSession - create and select the next numbered workout
//   deleteRecord - delete the active record while preserving ownership boundaries
// END_MODULE_MAP

package views

import (
	"context"
	"strconv"

	"charm.land/bubbles/v2/table"

	"github.com/danielsavinoff/trainings-logger/internal/app/command"
	"github.com/danielsavinoff/trainings-logger/internal/app/query"
)

// START_BLOCK_RECORDS_ACTIONS
func (m RecordsModel) selectionIDs() (string, string, string, string) {
	sessionID, catalogID, performanceID, setID := "", "", "", ""
	if i := m.selectedSession(); i >= 0 {
		sessionID = m.sessions[i].ID
	}

	if i := m.selectedCatalogExercise(); i >= 0 {
		catalogID = m.catalog[i].ID
	}

	if i := m.selectedPerformance(); i >= 0 {
		performanceID = m.performances[i].ID
	}

	if i := m.selectedSet(); i >= 0 {
		setID = m.sets[i].ID
	}

	return sessionID, catalogID, performanceID, setID
}

func (m *RecordsModel) reloadRecords(ctx context.Context, sessionID string, catalogID string, performanceID string, setID string) error {
	if m.store == nil {
		return nil
	}
	catalog, err := m.store.ListExercises.Handle(
		ctx,
		query.ListExercises{},
	)
	if err != nil {
		return err
	}
	sessions, err := m.store.ListSessions.Handle(
		ctx,
		query.ListSessions{},
	)
	if err != nil {
		return err
	}

	fallback := m.catalogTable.Cursor()
	m.catalog = catalog.Result
	m.catalogTable.SetRows(catalogRows(m.catalog))
	ids := make(
		[]string,
		0,
		len(m.catalog),
	)
	for _, item := range m.catalog {
		ids = append(
			ids,
			item.ID,
		)
	}
	restoreCursor(
		&m.catalogTable,
		ids,
		catalogID,
		fallback,
	)
	fallback = m.sessionTable.Cursor()
	m.sessions = sessions.Result
	rows := make(
		[]table.Row,
		0,
		len(m.sessions),
	)
	ids = nil
	for _, item := range m.sessions {
		rows = append(
			rows,
			table.Row{strconv.Itoa(item.Index)},
		)
		ids = append(
			ids,
			item.ID,
		)
	}
	m.sessionTable.SetRows(rows)
	restoreCursor(
		&m.sessionTable,
		ids,
		sessionID,
		fallback,
	)
	return m.loadPerformances(
		ctx,
		performanceID,
		setID,
	)
}

func (m *RecordsModel) addCatalogExercise() {
	if m.store == nil || m.transfer != nil {
		return
	}
	sessionID, catalogID, _, setID := m.selectionIDs()
	if sessionID == "" {
		m.err = "Create or select a session before adding an exercise."
		*m = m.resize()
		return
	}

	if catalogID == "" {
		return
	}
	result, err := m.store.AddPerformance.Handle(
		context.Background(),
		command.AddPerformance{
			SessionID:  sessionID,
			ExerciseID: catalogID,
		},
	)
	if err == nil {
		err = m.reloadRecords(
			context.Background(),
			sessionID,
			catalogID,
			result.Result.ID,
			setID,
		)
	}

	if err != nil {
		m.err = err.Error()
	} else {
		m.err = ""
		m.closeCatalog()
	}
	*m = m.resize()
}

func (m *RecordsModel) createSession() {
	if m.store == nil || m.transfer != nil {
		return
	}
	_, catalogID, _, _ := m.selectionIDs()
	result, err := m.store.CreateSession.Handle(
		context.Background(),
		command.CreateSession{},
	)
	if err == nil {
		err = m.reloadRecords(
			context.Background(),
			result.Result.ID,
			catalogID,
			"",
			"",
		)
	}

	if err != nil {
		m.err = err.Error()
	} else {
		m.err = ""
	}
	m.goPending = false
	*m = m.resize()
}

func (m *RecordsModel) deleteRecord() {
	if m.store == nil {
		return
	}
	sessionID, catalogID, performanceID, setID := m.selectionIDs()
	var err error
	switch m.activePanel {
	case sessionsPanel:
		if sessionID == "" {
			return
		}

		err = m.store.DeleteSession.Handle(
			context.Background(),
			command.DeleteSession{ID: sessionID},
		)
	case catalogPanel:
		if catalogID == "" {
			return
		}

		err = m.store.DeleteExercise.Handle(
			context.Background(),
			command.DeleteExercise{ID: catalogID},
		)
	case performancesPanel:
		if performanceID == "" {
			return
		}

		err = m.store.DeletePerformance.Handle(
			context.Background(),
			command.DeletePerformance{ID: performanceID},
		)
	case setsPanel:
		if setID == "" {
			return
		}

		err = m.store.DeleteSet.Handle(
			context.Background(),
			command.DeleteSet{ID: setID},
		)
	}

	if err == nil {
		err = m.reloadRecords(
			context.Background(),
			sessionID,
			catalogID,
			performanceID,
			setID,
		)
	}

	if err != nil {
		m.err = err.Error()
	} else {
		m.err = ""
	}
	*m = m.resize()
}

// END_BLOCK_RECORDS_ACTIONS
