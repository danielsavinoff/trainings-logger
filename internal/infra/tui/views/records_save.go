// START_MODULE_CONTRACT
//   PURPOSE: Validate and submit active-panel forms.
//   SCOPE: Invoke application commands and reload retained workout/catalog/performance/set IDs.
//   DEPENDS: M-APP-SESSIONS, M-APP-PERFORMANCES, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-EXERCISE, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   nonnegativeNumber - parse a finite nonnegative number
//   saveDialog - save a catalog definition or set
// END_MODULE_MAP

package views

import (
	"context"
	"math"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/app/command"
)

// START_BLOCK_RECORDS_SAVE
func nonnegativeNumber(value string) (float64, bool) {
	n, err := strconv.ParseFloat(
		value,
		64,
	)
	return n, err == nil && !math.IsNaN(n) && !math.IsInf(
		n,
		0,
	) && n >= 0
}

func (m *RecordsModel) saveDialog() tea.Cmd {
	d := m.dialog
	if m.store == nil {
		return d.Invalid(
			0,
			"Records storage is unavailable.",
		)
	}

	value := func(i int) string { return strings.TrimSpace(d.Value(i)) }
	sessionID, catalogID, performanceID, setID := m.selectionIDs()
	ctx := context.Background()
	var err error
	switch d.panel {
	case catalogPanel:
		name := value(0)
		if name == "" {
			return d.Invalid(
				0,
				"Enter an exercise name.",
			)
		}

		share := 0.0
		if value(1) != "" {
			var ok bool
			share, ok = nonnegativeNumber(value(1))
			if !ok || share > 100 {
				return d.Invalid(
					1,
					"Body weight share must be between 0 and 100%.",
				)
			}
		}

		if d.editing {
			err = m.store.UpdateExercise.Handle(
				ctx,
				command.UpdateExercise{
					ID:              d.recordID,
					Name:            name,
					BodyWeightShare: float32(share / 100),
				},
			)
			catalogID = d.recordID
		} else {
			result, saveErr := m.store.CreateExercise.Handle(
				ctx,
				command.CreateExercise{
					Name:            name,
					BodyWeightShare: float32(share / 100),
				},
			)
			err = saveErr
			catalogID = result.Result.ID
		}
	case setsPanel:
		number := func(i int) string {
			if value(i) == "" {
				return "0"
			}

			return value(i)
		}
		reps, parseErr := strconv.Atoi(number(0))
		if parseErr != nil || reps < 0 || reps > 255 {
			return d.Invalid(
				0,
				"Reps must be a whole number from 0 to 255.",
			)
		}
		weight, ok := nonnegativeNumber(number(1))
		if !ok {
			return d.Invalid(
				1,
				"Weight must be a number of 0 or more.",
			)
		}
		rest, parseErr := strconv.Atoi(number(2))
		if parseErr != nil || rest < 0 {
			return d.Invalid(
				2,
				"Rest must be a whole number of seconds, 0 or more.",
			)
		}

		if d.editing {
			err = m.store.UpdateSet.Handle(
				ctx,
				command.UpdateSet{
					ID:             d.recordID,
					Repetitions:    uint8(reps),
					ExerciseWeight: weight,
					WeightUnit:     m.weightUnit,
					RestInSeconds:  rest,
				},
			)
		} else {
			result, saveErr := m.store.CreateSet.Handle(
				ctx,
				command.CreateSet{
					PerformanceID:  d.performanceID,
					Repetitions:    uint8(reps),
					ExerciseWeight: weight,
					RestInSeconds:  rest,
				},
			)
			err = saveErr
			setID = result.Result.ID
		}
	}

	if err != nil {
		return d.Invalid(
			0,
			err.Error(),
		)
	}

	m.dialog = nil
	m.err = ""
	if err := m.reloadRecords(
		ctx,
		sessionID,
		catalogID,
		performanceID,
		setID,
	); err != nil {
		m.err = err.Error()
	}

	*m = m.resize()

	return nil
}

// END_BLOCK_RECORDS_SAVE
