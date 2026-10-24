package sqlite

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

type SetRepository struct {
	db *sql.DB
}

func NewSetRepository(db *sql.DB) *SetRepository {
	return &SetRepository{db: db}
}

func (r *SetRepository) Save(
	ctx context.Context,
	s set.Set,
) error {
	grams, err := s.ExerciseWeight.Weight.Grams()
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO sets (
			id, performance_id, exercise_weight_in_grams,
			repetitions, rest_in_seconds, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			exercise_weight_in_grams = excluded.exercise_weight_in_grams,
			repetitions = excluded.repetitions,
			rest_in_seconds = excluded.rest_in_seconds
	`

	_, err = persistence.ExecutorFromContext(
		ctx,
		r.db,
	).ExecContext(
		ctx,
		query,
		s.ID,
		s.PerformanceID,
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
		SELECT id, performance_id, exercise_weight_in_grams,
			repetitions, rest_in_seconds, created_at
		FROM sets
		WHERE id = ?
	`

	var storedID string
	var performanceID string
	var repetitions uint8
	var restInSeconds int
	var createdAt time.Time
	var grams int64
	err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&storedID,
		&performanceID,
		&grams,
		&repetitions,
		&restInSeconds,
		&createdAt,
	)
	if err != nil {
		return set.Set{}, err
	}

	return set.NewSet(
		storedID,
		performanceID,
		repetitions,
		float64(grams),
		shared.WeightUnitG,
		restInSeconds,
		&createdAt,
	)
}

func (r *SetRepository) FindAllByPerformanceID(
	ctx context.Context,
	performanceID string,
) ([]set.Set, error) {
	const query = `
		SELECT id, performance_id, exercise_weight_in_grams,
			repetitions, rest_in_seconds, created_at
		FROM sets WHERE performance_id = ?
		ORDER BY created_at DESC, id DESC
	`

	rows, err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryContext(
		ctx,
		query,
		performanceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []set.Set
	for rows.Next() {
		var id string
		var storedPerformanceID string
		var repetitions uint8
		var restInSeconds int
		var createdAt time.Time
		var grams int64
		err := rows.Scan(
			&id,
			&storedPerformanceID,
			&grams,
			&repetitions,
			&restInSeconds,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		s, err := set.NewSet(
			id,
			storedPerformanceID,
			repetitions,
			float64(grams),
			shared.WeightUnitG,
			restInSeconds,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(
			result,
			s,
		)
	}

	return result, rows.Err()
}

func (r *SetRepository) DeleteByID(
	ctx context.Context,
	id string,
) error {
	_, err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).ExecContext(
		ctx,
		"DELETE FROM sets WHERE id = ?",
		id,
	)

	return err
}
