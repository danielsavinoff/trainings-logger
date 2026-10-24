package command

import (
	"context"
	"database/sql"
	"errors"
	"uuid"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/performance"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
)

type CreateSet struct {
	PerformanceID  string
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
	performanceRepo performance.PerformanceRepository
}

func NewCreateSetHandler(
	preferencesRepo preferences.PreferencesRepository,
	setRepo set.SetRepository,
	performanceRepo performance.PerformanceRepository,
	transactions ports.TransactionRunner[CreateSetResult],
) CreateSetHandler {
	return CreateSetHandler{
		transactions:    transactions,
		preferencesRepo: preferencesRepo,
		setRepo:         setRepo,
		performanceRepo: performanceRepo,
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
				command.PerformanceID,
				command.Repetitions,
				command.ExerciseWeight,
				prefs.WeightUnit,
				command.RestInSeconds,
				nil,
			)
			if err != nil {
				return CreateSetResult{}, err
			}

			if _, err := h.performanceRepo.FindByID(
				ctx,
				command.PerformanceID,
			); err != nil {
				if errors.Is(
					err,
					sql.ErrNoRows,
				) {
					return CreateSetResult{}, set.ErrPerformanceOfSetNotFound
				}

				return CreateSetResult{}, err
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
