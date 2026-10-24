package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise"
)

type UpdateExercise struct {
	ID              string
	Name            string
	BodyWeightShare float32
}

type UpdateExerciseHandler struct {
	exerciseRepo exercise.ExerciseRepository
	transactions ports.TransactionRunner[ports.TransactionNoResult]
}

func NewUpdateExerciseHandler(
	exerciseRepo exercise.ExerciseRepository,
	transactions ports.TransactionRunner[ports.TransactionNoResult],
) UpdateExerciseHandler {
	return UpdateExerciseHandler{
		exerciseRepo: exerciseRepo,
		transactions: transactions,
	}
}

func (h UpdateExerciseHandler) Handle(
	ctx context.Context,
	command UpdateExercise,
) error {
	validated, err := exercise.NewExercise(
		command.ID,
		command.Name,
		command.BodyWeightShare,
	)
	if err != nil {
		return err
	}

	_, err = h.transactions.Run(
		ctx,
		func(
			ctx context.Context,
		) (ports.TransactionNoResult, error) {
			exists, err := h.exerciseRepo.ExistsByID(
				ctx,
				command.ID,
			)
			if err != nil {
				return ports.TransactionNoResult{}, err
			}

			if !exists {
				return ports.TransactionNoResult{}, exercise.ErrExerciseNotFound
			}

			return ports.TransactionNoResult{}, h.exerciseRepo.Save(
				ctx,
				validated,
			)
		},
	)

	return err
}
