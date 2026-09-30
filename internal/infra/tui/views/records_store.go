package views

import (
	"context"
	"math"
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/app/command"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"
)

type recordStore struct {
	tui.App
}

func (s *recordStore) saveExercise(
	record exerciseRecord,
	editing bool,
) (exerciseRecord, error) {
	e := exercise.Exercise{
		ID:              record.id,
		Name:            record.name,
		BodyWeightShare: record.bodyWeightShare,
	}

	if editing {
		err := s.UpdateExercise.Handle(
			context.Background(),
			command.UpdateExercise{
				ID:              e.ID,
				Name:            e.Name,
				BodyWeightShare: e.BodyWeightShare,
			},
		)
		if err != nil {
			return record, err
		}
	} else {
		result, err := s.CreateExercise.Handle(
			context.Background(),
			command.CreateExercise{
				Name:            e.Name,
				BodyWeightShare: e.BodyWeightShare,
			},
		)
		if err != nil {
			return record, err
		}

		e = result.Result
	}

	record.id = e.ID
	record.name = e.Name
	record.bodyWeightShare = e.BodyWeightShare

	return record, nil
}

func (s *recordStore) saveSet(
	exerciseID string,
	record workoutSet,
	editing bool,
	unit shared.WeightUnit,
) (workoutSet, error) {
	var saved set.Set

	if editing {
		var err error
		saved, err = set.NewSet(
			record.id,
			exerciseID,
			uint8(record.reps),
			record.weight,
			unit,
			int(record.rest/time.Second),
			&record.createdAt,
		)
		if err != nil {
			return record, err
		}

		err = s.UpdateSet.Handle(context.Background(), command.UpdateSet{
			ID:             saved.ID,
			ExerciseID:     saved.ExerciseID,
			Repetitions:    saved.Repetitions,
			ExerciseWeight: saved.ExerciseWeight.Weight.Value,
			WeightUnit:     saved.ExerciseWeight.Weight.Unit,
			RestInSeconds:  saved.RestInSeconds,
		})
		if err != nil {
			return record, err
		}
	} else {
		result, err := s.CreateSet.Handle(
			context.Background(),
			command.CreateSet{
				ExerciseID:     exerciseID,
				Repetitions:    uint8(record.reps),
				ExerciseWeight: record.weight,
				RestInSeconds:  int(record.rest / time.Second),
			},
		)
		if err != nil {
			return record, err
		}

		saved = result.Result
	}

	// Match the integer-gram precision of the stored value immediately.
	grams, err := saved.ExerciseWeight.Weight.Grams()
	if err != nil {
		return record, err
	}

	saved.ExerciseWeight.Weight = shared.Weight{
		Value: math.Round(grams),
		Unit:  shared.WeightUnitG,
	}

	return setRecord(saved, unit)
}

func setRecord(s set.Set, unit shared.WeightUnit) (workoutSet, error) {
	weight, err := s.ExerciseWeight.Weight.To(unit)
	if err != nil {
		return workoutSet{}, err
	}

	return workoutSet{
		id:        s.ID,
		createdAt: s.CreatedAt,
		reps:      int(s.Repetitions),
		weight:    weight.Value,
		rest:      time.Duration(s.RestInSeconds) * time.Second,
	}, nil
}
