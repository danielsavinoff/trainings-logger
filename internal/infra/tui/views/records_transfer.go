// START_MODULE_CONTRACT
//   PURPOSE: Save and load Dropbox backups from the records screen.
//   SCOPE: Asynchronous progress, operation ownership, records refresh, and selection retention.
//   DEPENDS: M-APP-REMOTE-STATE
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   recordTransfer - retain pending operation ownership
//   recordTransferMsg - return one completed transfer
//   startTransfer - start a bounded save/load without overlapping work
//   readTransferredRecords - reload records without startup setup
//   finishTransfer - show outcomes and retain selected IDs after a restore
// END_MODULE_MAP

package views

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"
)

// START_BLOCK_RECORDS_TRANSFER
type recordTransfer struct{ load bool }
type recordTransferMsg struct {
	owner   *recordTransfer
	records *RecordsModel
	err     error
}

func (m *RecordsModel) startTransfer(load bool) tea.Cmd {
	if m.transfer != nil {
		return nil
	}
	if m.store == nil || (load && m.store.LoadRemoteState == nil) || (!load && m.store.SaveRemoteState == nil) {
		m.transferStatus = "Dropbox transfer is not configured; press r to set up Dropbox."
		*m = m.resize()
		return nil
	}
	operation := &recordTransfer{load: load}
	m.transfer = operation
	m.goPending = false
	m.transferStatus = "Saving to Dropbox…"
	if load {
		m.transferStatus = "Loading from Dropbox…"
	}
	*m = m.resize()
	app := m.store.App
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			2*time.Minute,
		)
		defer cancel()
		if !load {
			return recordTransferMsg{owner: operation, err: app.SaveRemoteState(ctx)}
		}
		if err := app.LoadRemoteState(ctx); err != nil {
			return recordTransferMsg{owner: operation, err: err}
		}
		records, err := readTransferredRecords(
			ctx,
			app,
		)
		return recordTransferMsg{owner: operation, records: &records, err: err}
	}
}

func readTransferredRecords(
	ctx context.Context,
	app tui.App,
) (RecordsModel, error) {
	fresh, err := NewRecordsModel(
		ctx,
		app,
	)
	if err != nil {
		return RecordsModel{}, err
	}
	if fresh.questionnaire != nil {
		fresh.closeQuestionnaire()
	}
	for index := range fresh.exercises {
		fresh.table.SetCursor(index)
		if err := fresh.loadSelectedSets(ctx); err != nil {
			return RecordsModel{}, err
		}
	}
	fresh.table.SetCursor(0)
	return fresh, nil
}

func (m RecordsModel) finishTransfer(result recordTransferMsg) RecordsModel {
	if result.owner != m.transfer || m.transfer == nil {
		return m
	}
	operation := m.transfer
	m.transfer = nil
	if result.err != nil {
		m.transferStatus = result.err.Error()
		return m.resize()
	}
	if !operation.load {
		m.transferStatus = "Saved to Dropbox."
		return m.resize()
	}
	if result.records == nil {
		m.transferStatus = "Backup loaded; press q and reopen to refresh records."
		return m.resize()
	}
	exerciseID, setID := "", ""
	if index := m.selectedExercise(); index >= 0 {
		exerciseID = m.exercises[index].id
	}
	if session, index, ok := m.selectedSet(); ok {
		setID = m.exercises[m.selectedExercise()].sessions[session].sets[index].id
	}
	fresh := *result.records
	fresh.width, fresh.height = m.width, m.height
	fresh.showingHelp, fresh.help = m.showingHelp, m.help
	fresh.err = m.err
	fresh.transferStatus = "Loaded from Dropbox."
	for index, item := range fresh.exercises {
		if item.id == exerciseID {
			fresh.table.SetCursor(index)
			break
		}
	}
	if index := fresh.selectedExercise(); index >= 0 {
		fresh.setsExerciseID = fresh.exercises[index].id
		cursor := 0
		for _, session := range fresh.exercises[index].sessions {
			for _, item := range session.sets {
				if setID != "" && item.id == setID {
					fresh.setCursor = cursor
				}
				cursor++
			}
		}
	}
	fresh.focusPanel(m.activePanel)
	fresh = fresh.resize()
	fresh.syncExercisePager()
	fresh.refreshTimeline()
	return fresh
}

// END_BLOCK_RECORDS_TRANSFER
