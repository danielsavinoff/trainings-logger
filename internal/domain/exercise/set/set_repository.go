package set

import "context"

type SetRepository interface {
	Save(ctx context.Context, set Set) error
	FindByID(ctx context.Context, id string) (Set, error)
	FindAllByExerciseId(ctx context.Context, exerciseId string) ([]Set, error)
	DeleteByID(ctx context.Context, id string) error
}
