package exercise

import "context"

type ExerciseRepository interface {
	Save(ctx context.Context, exercise Exercise) error
	FindAll(ctx context.Context) ([]Exercise, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	DeleteByID(ctx context.Context, id string) error
}
