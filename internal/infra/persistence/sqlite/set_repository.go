package sqlite

import (
	"context"
	"database/sql"
	"math"

	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

type SetRepository struct {
	db *sql.DB
}

func NewSetRepository(db *sql.DB) *SetRepository {
	return &SetRepository{db: db}
}

func (r *SetRepository) Save(ctx context.Context, s set.Set) error {
	grams, err := s.ExerciseWeight.Weight.Grams()
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO sets (
			id, exercise_id, exercise_weight_in_grams,
			repetitions, rest_in_seconds, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			exercise_weight_in_grams = excluded.exercise_weight_in_grams,
			repetitions = excluded.repetitions,
			rest_in_seconds = excluded.rest_in_seconds
	`

	_, err = persistence.ExecutorFromContext(ctx, r.db).ExecContext(
		ctx,
		query,
		s.ID,
		s.ExerciseID,
		int64(math.Round(grams)),
		s.Repetitions,
		s.RestInSeconds,
		s.CreatedAt.UTC(),
	)

	return err
}

func (r *SetRepository) FindByID(
	ctx context.Context,
	id string,
) (set.Set, error) {
	const query = `
		SELECT id, exercise_id, exercise_weight_in_grams,
			repetitions, rest_in_seconds, created_at
		FROM sets
		WHERE id = ?
	`

	var result set.Set
	var grams int64
	err := persistence.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.ExerciseID,
		&grams,
		&result.Repetitions,
		&result.RestInSeconds,
		&result.CreatedAt,
	)
	if err != nil {
		return set.Set{}, err
	}

	result.ExerciseWeight.Weight = shared.Weight{
		Value: float64(grams),
		Unit:  shared.WeightUnitG,
	}

	return result, nil
}

func (r *SetRepository) FindAllByExerciseId(
	ctx context.Context,
	exerciseID string,
) ([]set.Set, error) {
	const query = `
		SELECT id, exercise_id, exercise_weight_in_grams,
			repetitions, rest_in_seconds, created_at
		FROM sets WHERE exercise_id = ?
		ORDER BY created_at DESC, rowid DESC
	`

	rows, err := persistence.ExecutorFromContext(ctx, r.db).QueryContext(ctx, query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []set.Set
	for rows.Next() {
		var s set.Set
		var grams int64
		err := rows.Scan(
			&s.ID,
			&s.ExerciseID,
			&grams,
			&s.Repetitions,
			&s.RestInSeconds,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		s.ExerciseWeight.Weight = shared.Weight{
			Value: float64(grams),
			Unit:  shared.WeightUnitG,
		}
		result = append(result, s)
	}

	return result, rows.Err()
}

func (r *SetRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := persistence.ExecutorFromContext(ctx, r.db).ExecContext(
		ctx, "DELETE FROM sets WHERE id = ?", id,
	)

	return err
}
