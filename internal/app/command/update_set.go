package command

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
)

type UpdateSet struct {
	ID             string
	ExerciseID     string
	Repetitions    uint8
	ExerciseWeight float64
	WeightUnit     shared.WeightUnit
	RestInSeconds  int
}

type UpdateSetHandler struct {
	setRepo      set.SetRepository
	transactions ports.TransactionRunner[ports.TransactionNoResult]
}

func NewUpdateSetHandler(
	setRepo set.SetRepository,
	transactions ports.TransactionRunner[ports.TransactionNoResult],
) UpdateSetHandler {
	return UpdateSetHandler{
		setRepo:      setRepo,
		transactions: transactions,
	}
}

func (h UpdateSetHandler) Handle(
	ctx context.Context,
	command UpdateSet,
) error {
	_, err := h.transactions.Run(ctx, func(
		ctx context.Context,
	) (ports.TransactionNoResult, error) {
		existing, err := h.setRepo.FindByID(ctx, command.ID)
		if errors.Is(err, sql.ErrNoRows) {
			return ports.TransactionNoResult{}, set.ErrSetNotFound
		}
		if err != nil {
			return ports.TransactionNoResult{}, err
		}

		validated, err := set.NewSet(
			command.ID,
			existing.ExerciseID,
			command.Repetitions,
			command.ExerciseWeight,
			command.WeightUnit,
			command.RestInSeconds,
			nil,
		)
		if err != nil {
			return ports.TransactionNoResult{}, err
		}

		return ports.TransactionNoResult{}, h.setRepo.Save(ctx, validated)
	})

	return err
}
