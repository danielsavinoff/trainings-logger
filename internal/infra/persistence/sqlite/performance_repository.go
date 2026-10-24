// START_MODULE_CONTRACT
//   PURPOSE: Persist exercises added to workout sessions.
//   SCOPE: Constructor-validated workout exercise membership and transactional persistence.
//   DEPENDS: M-DOMAIN-PERFORMANCE
//   LINKS: M-INFRA-PERSISTENCE-SQLITE, V-M-INFRA-PERSISTENCE-SQLITE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   PerformanceRepository - persist unique workout exercise membership
//   performanceScanner - accept row and rows scanners
//   NewPerformanceRepository - wire the SQLite repository
//   scanPerformance - reconstruct a validated domain membership
//   Save - insert membership once per workout and catalog exercise
//   FindByID - read a validated performance
//   FindBySessionAndExercise - read a unique membership pair
//   FindAllBySessionID - list workout performances in insertion order
//   DeleteByID - remove one performance and cascade only its sets
// END_MODULE_MAP

package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/performance"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

// START_BLOCK_PERFORMANCE_REPOSITORY
type PerformanceRepository struct {
	db *sql.DB
}

type performanceScanner interface {
	Scan(...any) error
}

func NewPerformanceRepository(db *sql.DB) *PerformanceRepository {
	return &PerformanceRepository{db: db}
}

func scanPerformance(row performanceScanner) (performance.Performance, error) {
	var id, sessionID, exerciseID string
	var createdAt time.Time
	if err := row.Scan(
		&id,
		&sessionID,
		&exerciseID,
		&createdAt,
	); err != nil {
		return performance.Performance{}, err
	}

	return performance.NewPerformance(
		id,
		sessionID,
		exerciseID,
		&createdAt,
	)
}

func (r *PerformanceRepository) Save(ctx context.Context, item performance.Performance) error {
	validated, err := performance.NewPerformance(
		item.ID,
		item.SessionID,
		item.ExerciseID,
		&item.CreatedAt,
	)
	if err != nil {
		return err
	}
	_, err = persistence.ExecutorFromContext(
		ctx,
		r.db,
	).ExecContext(
		ctx,
		`INSERT INTO performances(id,session_id,exercise_id,created_at) VALUES (?,?,?,?)
  ON CONFLICT(session_id,exercise_id) DO NOTHING`,
		validated.ID,
		validated.SessionID,
		validated.ExerciseID,
		validated.CreatedAt,
	)
	return err
}

func (r *PerformanceRepository) FindByID(ctx context.Context, id string) (performance.Performance, error) {
	return scanPerformance(persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryRowContext(
		ctx,
		`SELECT id,session_id,exercise_id,created_at FROM performances WHERE id=?`,
		id,
	))
}

func (r *PerformanceRepository) FindBySessionAndExercise(ctx context.Context, sessionID string, exerciseID string) (performance.Performance, error) {
	return scanPerformance(persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryRowContext(
		ctx,
		`SELECT id,session_id,exercise_id,created_at FROM performances WHERE session_id=? AND exercise_id=?`,
		sessionID,
		exerciseID,
	))
}

func (r *PerformanceRepository) FindAllBySessionID(ctx context.Context, sessionID string) ([]performance.Performance, error) {
	rows, err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryContext(
		ctx,
		`SELECT id,session_id,exercise_id,created_at FROM performances WHERE session_id=? ORDER BY created_at ASC,id ASC`,
		sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []performance.Performance{}
	for rows.Next() {
		item, err := scanPerformance(rows)
		if err != nil {
			return nil, err
		}

		result = append(
			result,
			item,
		)
	}

	return result, rows.Err()
}

func (r *PerformanceRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).ExecContext(
		ctx,
		`DELETE FROM performances WHERE id=?`,
		id,
	)
	return err
}

var _ performance.PerformanceRepository = (*PerformanceRepository)(nil)

// END_BLOCK_PERFORMANCE_REPOSITORY
