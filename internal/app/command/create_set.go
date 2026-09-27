package command

import (
	"context"
	"uuid"

	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
)

type CreateSet struct {
	ExerciseID     string
	ExerciseWeight float64
	Repetitions    uint8
}

type CreateSetResult struct {
	Result set.Set
}

type CreateSetHandler struct {
	preferencesRepo preferences.PreferencesRepository
	setRepo         set.SetRepository
	exerciseRepo    exercise.ExerciseRepository
}

func NewCreateSetHandler(
	preferencesRepo preferences.PreferencesRepository,
	setRepo set.SetRepository,
	exerciseRepo exercise.ExerciseRepository,
) CreateSetHandler {
	return CreateSetHandler{
		preferencesRepo: preferencesRepo,
		setRepo:         setRepo,
		exerciseRepo:    exerciseRepo,
	}
}

func (h CreateSetHandler) Handle(ctx context.Context, command CreateSet) (CreateSetResult, error) {
	preferences, err := h.preferencesRepo.Find(ctx)
	if err != nil {
		return CreateSetResult{}, err
	}

	setID := uuid.NewV4()
	userSet, err := set.NewSet(
		setID.String(),
		command.Repetitions,
		command.ExerciseWeight,
		preferences.WeightUnit,
		nil,
	)
	if err != nil {
		return CreateSetResult{}, err
	}

	exists, err := h.exerciseRepo.ExistsByID(ctx, command.ExerciseID)
	if err != nil {
		return CreateSetResult{}, err
	}
	if exists != true {
		return CreateSetResult{}, set.ErrExerciseOfSetNotFound
	}

	err = h.setRepo.Save(ctx, userSet)
	if err != nil {
		return CreateSetResult{}, err
	}

	return CreateSetResult{
		Result: userSet,
	}, nil
}
