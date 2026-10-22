package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
)

type UpdatePreferences struct {
	WeightUnit shared.WeightUnit
}

type UpdatePreferencesHandler struct {
	repository   preferences.PreferencesRepository
	transactions ports.TransactionRunner[ports.TransactionNoResult]
}

func NewUpdatePreferencesHandler(
	repository preferences.PreferencesRepository,
	transactions ports.TransactionRunner[ports.TransactionNoResult],
) UpdatePreferencesHandler {
	return UpdatePreferencesHandler{
		repository:   repository,
		transactions: transactions,
	}
}

func (h UpdatePreferencesHandler) Handle(
	ctx context.Context,
	command UpdatePreferences,
) error {
	validated, err := preferences.NewPreferences(command.WeightUnit)
	if err != nil {
		return err
	}

	_, err = h.transactions.Run(
		ctx,
		func(
			ctx context.Context,
		) (ports.TransactionNoResult, error) {
			return ports.TransactionNoResult{}, h.repository.Save(
				ctx,
				validated,
			)
		},
	)

	return err
}
