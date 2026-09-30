package command

import (
	"context"
	"uuid"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise"
)

type CreateExercise struct {
	Name            string
	BodyWeightShare float32
}

type CreateExerciseResult struct {
	Result exercise.Exercise
}

type CreateExerciseHandler struct {
	exerciseRepo exercise.ExerciseRepository
	transactions ports.TransactionRunner[CreateExerciseResult]
}

func NewCreateExerciseHandler(
	exerciseRepo exercise.ExerciseRepository,
	transactions ports.TransactionRunner[CreateExerciseResult],
) CreateExerciseHandler {
	return CreateExerciseHandler{
		exerciseRepo: exerciseRepo,
		transactions: transactions,
	}
}

func (h CreateExerciseHandler) Handle(
	ctx context.Context,
	command CreateExercise,
) (CreateExerciseResult, error) {
	exerciseId := uuid.NewV4()
	userExercise, err := exercise.NewExercise(
		exerciseId.String(),
		command.Name,
		command.BodyWeightShare,
	)
	if err != nil {
		return CreateExerciseResult{}, err
	}

	return h.transactions.Run(ctx, func(
		ctx context.Context,
	) (CreateExerciseResult, error) {
		exists, err := h.exerciseRepo.ExistsByName(ctx, userExercise.Name)
		if err != nil {
			return CreateExerciseResult{}, err
		}

		if exists {
			return CreateExerciseResult{}, exercise.ErrExerciseNameConflict
		}

		if err := h.exerciseRepo.Save(ctx, userExercise); err != nil {
			return CreateExerciseResult{}, err
		}

		return CreateExerciseResult{Result: userExercise}, nil
	})
}
