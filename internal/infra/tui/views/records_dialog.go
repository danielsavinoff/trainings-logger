package views

import (
	"strconv"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

type recordDialog struct {
	*components.Dialog
	panel    panel
	editing  bool
	exercise int
	session  int
	set      int
}

func (m *RecordsModel) openDialog(editing bool) tea.Cmd {
	exercise := m.selectedExercise()
	if (editing || m.activePanel == setsPanel) && exercise < 0 {
		return nil
	}

	d := &recordDialog{panel: m.activePanel, editing: editing, exercise: exercise}
	action := "New"
	var title string
	var fields []components.Field
	if editing {
		action = "Edit"
	}

	if d.panel == exercisesPanel {
		title = action + " exercise"
		name, bodyWeight := "", ""
		if exercise >= 0 {
			e := m.exercises[exercise]
			if editing {
				name = e.name
			}

			if editing {
				bodyWeight = strconv.FormatFloat(
					float64(e.bodyWeightShare)*100, 'g', 5, 64,
				)
			}
		}

		fields = []components.Field{
			components.NewField("Name", name, "Exercise name"),
			components.NewField("Body weight (%)", bodyWeight, "0–100 (bench: 0)"),
		}
	} else {
		title = action + " set"
		set := workoutSet{reps: 8}
		session, index, ok := m.selectedSet()
		if editing && !ok {
			return nil
		}

		if ok {
			d.session, d.set = session, index
			set = m.exercises[exercise].sessions[session].sets[index]
		}

		fields = []components.Field{
			components.NewField("Reps", strconv.Itoa(set.reps), "8"),
			components.NewField(
				"Weight ("+string(m.weightUnit)+")",
				strconv.FormatFloat(set.weight, 'f', -1, 64),
				"0",
			),
			components.NewField(
				"Rest (seconds)",
				strconv.Itoa(int(set.rest/time.Second)),
				"150",
			),
		}
	}

	m.goPending = false
	d.Dialog = components.NewDialog(title, fields)
	m.dialog = d
	d.Resize(m.width, m.height)

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
