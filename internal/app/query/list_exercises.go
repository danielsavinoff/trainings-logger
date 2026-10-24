package query

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise"
)

type ListExercises struct{}

type ListExercisesResult struct {
	Result []exercise.Exercise
}

type ListExercisesHandler struct {
	transactions ports.TransactionRunner[ListExercisesResult]
	exerciseRepo exercise.ExerciseRepository
}

func NewListExercisesHandler(
	repo exercise.ExerciseRepository,
	transactions ports.TransactionRunner[ListExercisesResult],
) ListExercisesHandler {
	return ListExercisesHandler{
		exerciseRepo: repo,
		transactions: transactions,
	}
}

func (h ListExercisesHandler) Handle(
	ctx context.Context,
	query ListExercises,
) (ListExercisesResult, error) {
	return h.transactions.Run(
		ctx,
		func(
			ctx context.Context,
		) (ListExercisesResult, error) {
			exercises, err := h.exerciseRepo.FindAll(ctx)
			if err != nil {
				return ListExercisesResult{}, err
			}

			return ListExercisesResult{Result: exercises}, nil
		},
	)
}
