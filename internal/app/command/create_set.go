package command

import (
	"context"
	"errors"
	"uuid"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
)

type CreateSet struct {
	ExerciseID     string
	ExerciseWeight float64
	RestInSeconds  int
	Repetitions    uint8
}

type CreateSetResult struct {
	Result set.Set
}

type CreateSetHandler struct {
	transactions    ports.TransactionRunner[CreateSetResult]
	preferencesRepo preferences.PreferencesRepository
	setRepo         set.SetRepository
	exerciseRepo    exercise.ExerciseRepository
}

func NewCreateSetHandler(
	preferencesRepo preferences.PreferencesRepository,
	setRepo set.SetRepository,
	exerciseRepo exercise.ExerciseRepository,
	transactions ports.TransactionRunner[CreateSetResult],
) CreateSetHandler {
	return CreateSetHandler{
		transactions:    transactions,
		preferencesRepo: preferencesRepo,
		setRepo:         setRepo,
		exerciseRepo:    exerciseRepo,
	}
}

func (h CreateSetHandler) Handle(
	ctx context.Context,
	command CreateSet,
) (CreateSetResult, error) {
	return h.transactions.Run(
		ctx,
		func(
			ctx context.Context,
		) (CreateSetResult, error) {
			prefs, err := h.preferencesRepo.Find(ctx)
			if errors.Is(
				err,
				preferences.ErrPreferencesNotFound,
			) {
				prefs, err = preferences.NewPreferences(shared.WeightUnitKg)
			}

			if err != nil {
				return CreateSetResult{}, err
			}

			userSet, err := set.NewSet(
				uuid.NewV4().String(),
				command.ExerciseID,
				command.Repetitions,
				command.ExerciseWeight,
				prefs.WeightUnit,
				command.RestInSeconds,
				nil,
			)
			if err != nil {
				return CreateSetResult{}, err
			}

			exists, err := h.exerciseRepo.ExistsByID(
				ctx,
				command.ExerciseID,
			)
			if err != nil {
				return CreateSetResult{}, err
			}

			if !exists {
				return CreateSetResult{}, set.ErrExerciseOfSetNotFound
			}

			if err := h.setRepo.Save(
				ctx,
				userSet,
			); err != nil {
				return CreateSetResult{}, err
			}

			return CreateSetResult{Result: userSet}, nil
		},
	)
}
