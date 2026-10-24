package performance

import "context"

type PerformanceRepository interface {
	Save(ctx context.Context, item Performance) error
	FindByID(ctx context.Context, id string) (Performance, error)
	FindBySessionAndExercise(ctx context.Context, sessionID string, exerciseID string) (Performance, error)
	FindAllBySessionID(ctx context.Context, sessionID string) ([]Performance, error)
	DeleteByID(ctx context.Context, id string) error
}
