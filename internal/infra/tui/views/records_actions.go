package views

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/command"
)

func (m *RecordsModel) refreshRecords(exercise int) {
	m.table.SetRows(exerciseRows(m.exercises))
	m.table.SetCursor(max(0, min(exercise, len(m.exercises)-1)))
	m.syncExercisePager()

	selectedID := ""
	if index := m.selectedExercise(); index >= 0 {
		selectedID = m.exercises[index].id
	}

	if selectedID != m.setsExerciseID {
		m.refreshSelectedSets()

		return
	}

	m.refreshTimeline()
}

func (m *RecordsModel) refreshSelectedSets() {
	if err := m.loadSelectedSets(context.Background()); err != nil {
		m.err = err.Error()
	} else {
		m.err = ""
	}

	m.refreshTimeline()
}

func (m *RecordsModel) deleteRecord() {
	exercise := m.selectedExercise()
	if exercise < 0 {
		return
	}

	if m.store != nil {
		var err error
		if m.activePanel == exercisesPanel {
			err = m.store.DeleteExercise.Handle(
				context.Background(),
				command.DeleteExercise{ID: m.exercises[exercise].id},
			)
		} else {
			session, index, ok := m.selectedSet()
			if !ok {
				return
			}

			err = m.store.DeleteSet.Handle(
				context.Background(),
				command.DeleteSet{
					ID: m.exercises[exercise].sessions[session].sets[index].id,
				},
			)
		}

		if err != nil {
			m.err = err.Error()

			return
		}
	}

	m.err = ""

	if m.activePanel == exercisesPanel {
		m.exercises = append(m.exercises[:exercise], m.exercises[exercise+1:]...)
		m.setCursor = 0
	} else {
		session, set, ok := m.selectedSet()
		if !ok {
			return
		}

		e := &m.exercises[exercise]
		s := &e.sessions[session]
		s.sets = append(s.sets[:set], s.sets[set+1:]...)
		if len(s.sets) == 0 {
			e.sessions = append(e.sessions[:session], e.sessions[session+1:]...)
		}
	}

	m.refreshRecords(exercise)
}
