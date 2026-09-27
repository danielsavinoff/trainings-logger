package command

import (
	"context"
	"uuid"

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
}

func NewCreateExerciseHandler(
	exerciseRepo exercise.ExerciseRepository,
) CreateExerciseHandler {
	return CreateExerciseHandler{
		exerciseRepo: exerciseRepo,
	}
}

func (h CreateExerciseHandler) Handle(ctx context.Context, command CreateExercise) (CreateExerciseResult, error) {
	exerciseId := uuid.NewV4()
	userExercise, err := exercise.NewExercise(
		exerciseId.String(),
		command.Name,
		command.BodyWeightShare,
	)
	if err != nil {
		return CreateExerciseResult{}, err
	}

	exists, err := h.exerciseRepo.ExistsByName(ctx, command.Name)
	if exists != false {
		return CreateExerciseResult{}, exercise.ErrExerciseNameConflict
	}

	err = h.exerciseRepo.Save(ctx, userExercise)
	if err != nil {
		return CreateExerciseResult{}, err
	}

	return CreateExerciseResult{
		Result: userExercise,
	}, nil
}
