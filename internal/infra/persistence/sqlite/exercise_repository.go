package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
	"modernc.org/sqlite"
)

type ExerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepository(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) Save(
	ctx context.Context,
	e exercise.Exercise,
) error {
	const query = `
		INSERT INTO exercises (id, name, body_weight_share)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			body_weight_share = excluded.body_weight_share
	`

	_, err := persistence.ExecutorFromContext(ctx, r.db).ExecContext(
		ctx, query, e.ID, e.Name, e.BodyWeightShare,
	)

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == 2067 {
		return exercise.ErrExerciseNameConflict
	}

	return err
}

func (r *ExerciseRepository) FindAll(
	ctx context.Context,
) ([]exercise.Exercise, error) {
	const query = `
		SELECT id, name, body_weight_share FROM exercises ORDER BY rowid
	`

	rows, err := persistence.ExecutorFromContext(ctx, r.db).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []exercise.Exercise
	for rows.Next() {
		var e exercise.Exercise
		if err := rows.Scan(&e.ID, &e.Name, &e.BodyWeightShare); err != nil {
			return nil, err
		}

		result = append(result, e)
	}

	return result, rows.Err()
}

func (r *ExerciseRepository) ExistsByName(
	ctx context.Context,
	name string,
) (bool, error) {
	const query = "SELECT EXISTS(SELECT 1 FROM exercises WHERE name = ?)"

	var exists bool
	err := persistence.ExecutorFromContext(ctx, r.db).
		QueryRowContext(ctx, query, name).
		Scan(&exists)

	return exists, err
}

func (r *ExerciseRepository) ExistsByID(
	ctx context.Context,
	id string,
) (bool, error) {
	const query = "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = ?)"

	var exists bool
	err := persistence.ExecutorFromContext(ctx, r.db).
		QueryRowContext(ctx, query, id).
		Scan(&exists)

	return exists, err
}

func (r *ExerciseRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := persistence.ExecutorFromContext(ctx, r.db).ExecContext(
		ctx, "DELETE FROM exercises WHERE id = ?", id,
	)

	return err
}
