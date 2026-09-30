package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise"
)

type DeleteExercise struct {
	ID string
}

type DeleteExerciseHandler struct {
	exerciseRepo exercise.ExerciseRepository
	transactions ports.TransactionRunner[ports.TransactionNoResult]
}

func NewDeleteExerciseHandler(
	exerciseRepo exercise.ExerciseRepository,
	transactions ports.TransactionRunner[ports.TransactionNoResult],
) DeleteExerciseHandler {
	return DeleteExerciseHandler{
		exerciseRepo: exerciseRepo,
		transactions: transactions,
	}
}

func (h DeleteExerciseHandler) Handle(
	ctx context.Context,
	command DeleteExercise,
) error {
	_, err := h.transactions.Run(ctx, func(
		ctx context.Context,
	) (ports.TransactionNoResult, error) {
		return ports.TransactionNoResult{},
			h.exerciseRepo.DeleteByID(ctx, command.ID)
	})

	return err
}
