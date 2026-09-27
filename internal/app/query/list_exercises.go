package query

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise"
)

type ListExercises struct{}

type ListExercisesResult struct {
	Result []exercise.Exercise
}

type ListExercisesHandler struct {
	exerciseRepo exercise.ExerciseRepository
}

func (h ListExercisesHandler) Handle(ctx context.Context, query ListExercises) (ListExercisesResult, error) {
	exercises, err := h.exerciseRepo.FindAll(ctx)
	if err != nil {
		return ListExercisesResult{}, err
	}

	return ListExercisesResult{
		Result: exercises,
	}, nil
}
