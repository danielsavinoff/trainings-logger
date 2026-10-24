// START_MODULE_CONTRACT
//   PURPOSE: Open the catalog selector or catalog and set forms.
//   SCOPE: Automatic session creation, shared catalog edits, and set edits under a captured performance.
//   DEPENDS: M-APP-SESSIONS, M-APP-PERFORMANCES, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-EXERCISE, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   recordDialog - retain generic form state and performance ownership
//   openDialog - open catalog selection or construct the active-panel form
//   updateDialog - route form submission or cancellation
// END_MODULE_MAP

package views

import (
	"strconv"

	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_DIALOG
type recordDialog struct {
	*components.Dialog
	panel         panel
	editing       bool
	recordID      string
	performanceID string
}

func (m *RecordsModel) openDialog(editing bool) tea.Cmd {
	if m.activePanel == sessionsPanel {
		if !editing {
			m.createSession()
		}

		return nil
	}

	if m.activePanel == performancesPanel {
		if !editing {
			m.openCatalog()
		}

		return nil
	}
	_, _, performanceID, setID := m.selectionIDs()
	if m.activePanel == setsPanel && performanceID == "" {
		m.err = "Add a performance to the session before recording sets."
		*m = m.resize()
		return nil
	}

	d := &recordDialog{
		panel:         m.activePanel,
		editing:       editing,
		performanceID: performanceID,
	}

	action := "New"
	if editing {
		action = "Edit"
	}

	var title string
	var fields []components.Field
	switch m.activePanel {
	case catalogPanel:
		name, share := "", ""
		if editing {
			i := m.selectedCatalogExercise()
			if i < 0 {
				return nil
			}

			item := m.catalog[i]
			d.recordID = item.ID
			name = item.Name
			share = strconv.FormatFloat(
				float64(item.BodyWeightShare)*100,
				'g',
				5,
				64,
			)
		}

		title = action + " catalog exercise"
		fields = []components.Field{components.NewField(
			"Name",
			name,
			"Exercise name",
		), components.NewField(
			"Body weight (%)",
			share,
			"0–100 (bench: 0)",
		)}
	case setsPanel:
		if editing && setID == "" {
			return nil
		}
		reps, weight, rest := "8", "0", "0"
		if i := m.selectedSet(); i >= 0 {
			item := m.sets[i]
			display, err := item.ExerciseWeight.Weight.To(m.weightUnit)
			if err != nil {
				m.err = err.Error()
				*m = m.resize()
				return nil
			}

			reps = strconv.Itoa(int(item.Repetitions))
			weight = strconv.FormatFloat(
				display.Value,
				'f',
				-1,
				64,
			)
			rest = strconv.Itoa(item.RestInSeconds)
			if editing {
				d.recordID = item.ID
			}
		}

		title = action + " set"
		fields = []components.Field{components.NewField(
			"Reps",
			reps,
			"8",
		), components.NewField(
			"Weight ("+string(m.weightUnit)+")",
			weight,
			"0",
		), components.NewField(
			"Rest (seconds)",
			rest,
			"150",
		)}
	}

	m.goPending = false
	d.Dialog = components.NewDialog(
		title,
		fields,
	)
	m.dialog = d
	d.Resize(
		m.width,
		m.height,
	)
	return d.Init()
}

func (m RecordsModel) updateDialog(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.dialog.Dialog, cmd = m.dialog.Dialog.Update(msg)
	switch m.dialog.State {
	case components.DialogClosed:
		m.dialog = nil
	case components.DialogSubmitted:
		cmd = m.saveDialog()
	}

	return m, cmd
}

// END_BLOCK_RECORDS_DIALOG
