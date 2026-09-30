package query

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
)

type GetPreferences struct{}

type GetPreferencesResult struct {
	Result preferences.Preferences
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
	return h.transactions.Run(ctx, func(
		ctx context.Context,
	) (GetPreferencesResult, error) {
		prefs, err := h.preferencesRepo.Find(ctx)
		if err != nil {
			return GetPreferencesResult{}, err
		}

		return GetPreferencesResult{Result: prefs}, nil
	})
}
