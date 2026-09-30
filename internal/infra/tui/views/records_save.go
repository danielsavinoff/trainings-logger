package views

import (
	"math"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

func nonnegativeNumber(value string) (float64, bool) {
	n, err := strconv.ParseFloat(value, 64)

	return n, err == nil && !math.IsNaN(n) && !math.IsInf(n, 0) && n >= 0
}

func (m *RecordsModel) saveDialog() tea.Cmd {
	d := m.dialog
	value := func(i int) string {
		return strings.TrimSpace(d.Value(i))
	}

	if d.panel == exercisesPanel {
		name := value(0)
		if name == "" {
			return d.Invalid(0, "Enter an exercise name.")
		}

		bodyWeight := 0.0
		if value(1) != "" {
			var ok bool
			bodyWeight, ok = nonnegativeNumber(value(1))
			if !ok || bodyWeight > 100 {
				return d.Invalid(1, "Body weight share must be between 0 and 100%.")
			}
		}

		record := exerciseRecord{
			name:            name,
			bodyWeightShare: float32(bodyWeight / 100),
		}

		if d.editing {
			record.id = m.exercises[d.exercise].id
			record.sessions = m.exercises[d.exercise].sessions
		}

		if m.store != nil {
			saved, err := m.store.saveExercise(record, d.editing)
			if err != nil {
				return d.Invalid(0, err.Error())
			}

			record = saved
		}

		if d.editing {
			m.exercises[d.exercise] = record
		} else {
			m.exercises = append(m.exercises, record)
			d.exercise = len(m.exercises) - 1
			m.setCursor = 0
		}
	} else {
		number := func(i int) string {
			if value(i) == "" {
				return "0"
			}

			return value(i)
		}

		reps, err := strconv.Atoi(number(0))
		if err != nil || reps < 0 || reps > 255 {
			return d.Invalid(0, "Reps must be a whole number from 0 to 255.")
		}

		weight, ok := nonnegativeNumber(number(1))
		if !ok {
			return d.Invalid(1, "Weight must be a number of 0 or more.")
		}

		seconds, err := strconv.ParseInt(number(2), 10, 64)
		if err != nil || seconds < 0 ||
			seconds > int64(math.MaxInt64)/int64(time.Second) {
			return d.Invalid(2, "Rest must be a whole number of seconds, 0 or more.")
		}

		set := workoutSet{
			reps:      reps,
			weight:    weight,
			rest:      time.Duration(seconds) * time.Second,
			createdAt: time.Now(),
		}

		if d.editing {
			original := m.exercises[d.exercise].sessions[d.session].sets[d.set]
			set.id, set.createdAt = original.id, original.createdAt
		}

		if m.store != nil {
			saved, err := m.store.saveSet(
				m.exercises[d.exercise].id, set, d.editing, m.weightUnit,
			)
			if err != nil {
				return d.Invalid(0, err.Error())
			}

			set = saved
		}

		if d.editing {
			m.exercises[d.exercise].sessions[d.session].sets[d.set] = set
		} else {
			m.appendSet(d.exercise, set, set.createdAt.Local())
		}
	}

	m.err = ""
	m.refreshRecords(d.exercise)
	m.dialog = nil

	return nil
}
