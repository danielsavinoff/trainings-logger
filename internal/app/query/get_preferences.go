package query

import (
	"context"
	"errors"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
)

type GetPreferences struct{}

type GetPreferencesResult struct {
	Result preferences.Preferences
	Exists bool
}

type GetPreferencesHandler struct {
	preferencesRepo preferences.PreferencesRepository
	transactions    ports.TransactionRunner[GetPreferencesResult]
}

func NewGetPreferencesHandler(
	preferencesRepo preferences.PreferencesRepository,
	transactions ports.TransactionRunner[GetPreferencesResult],
) GetPreferencesHandler {
	return GetPreferencesHandler{
		preferencesRepo: preferencesRepo,
		transactions:    transactions,
	}
}

func (h GetPreferencesHandler) Handle(
	ctx context.Context,
	query GetPreferences,
) (GetPreferencesResult, error) {
	return h.transactions.Run(
		ctx,
		func(
			ctx context.Context,
		) (GetPreferencesResult, error) {
			prefs, err := h.preferencesRepo.Find(ctx)
			if errors.Is(
				err,
				preferences.ErrPreferencesNotFound,
			) {
				defaults, err := preferences.NewPreferences(shared.WeightUnitKg)
				if err != nil {
					return GetPreferencesResult{}, err
				}

				return GetPreferencesResult{Result: defaults, Exists: false}, nil
			}

			if err != nil {
				return GetPreferencesResult{}, err
			}

			return GetPreferencesResult{Result: prefs, Exists: true}, nil
		},
	)
}
