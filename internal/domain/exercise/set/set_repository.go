package set

import "context"

type SetRepository interface {
	Save(ctx context.Context, set Set) error
	FindAllByExerciseId(exerciseId string) ([]Set, error)
}
