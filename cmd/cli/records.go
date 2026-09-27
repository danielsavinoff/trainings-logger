package main

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/table"
)

func exerciseRows(exercises []mockExercise) []table.Row {
	rows := make([]table.Row, 0, len(exercises))
	for _, exercise := range exercises {
		percentage := "—"
		if len(exercise.sessions) > 0 && len(exercise.sessions[0].sets) > 0 && exercise.bodyWeightKG > 0 {
			percentage = fmt.Sprintf("%.1f%%", exercise.sessions[0].sets[0].weightKG/exercise.bodyWeightKG*100)
		}
		rows = append(rows, table.Row{exercise.name, percentage})
	}
	return rows
}

func (m model) selectedExercise() int {
	i := m.table.Cursor()
	if i < 0 || i >= len(m.exercises) {
		return -1
	}
	return i
}

func (m model) selectedSet() (session, set int, ok bool) {
	exercise := m.selectedExercise()
	if exercise < 0 {
		return 0, 0, false
	}
	index := m.setCursor
	for s, session := range m.exercises[exercise].sessions {
		if index >= 0 && index < len(session.sets) {
			return s, index, true
		}
		index -= len(session.sets)
	}
	return 0, 0, false
}

func (m *model) refreshRecords(exercise int) {
	m.table.SetRows(exerciseRows(m.exercises))
	m.table.SetCursor(max(0, min(exercise, len(m.exercises)-1)))
	m.syncExercisePager()
	m.refreshTimeline()
}

func (m *model) deleteRecord() {
	exercise := m.selectedExercise()
	if exercise < 0 {
		return
	}
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

func (m *model) appendSet(exercise int, set workoutSet, now time.Time) {
	e := &m.exercises[exercise]
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	index, cursor := 0, 0
	for index < len(e.sessions) && e.sessions[index].date.Format(time.DateOnly) > date.Format(time.DateOnly) {
		cursor += len(e.sessions[index].sets)
		index++
	}
	if index < len(e.sessions) && e.sessions[index].date.Format(time.DateOnly) == date.Format(time.DateOnly) {
		cursor += len(e.sessions[index].sets)
		e.sessions[index].sets = append(e.sessions[index].sets, set)
	} else {
		e.sessions = append(e.sessions, workoutSession{})
		copy(e.sessions[index+1:], e.sessions[index:])
		e.sessions[index] = workoutSession{date: date, sets: []workoutSet{set}}
	}
	m.setCursor = cursor
}
