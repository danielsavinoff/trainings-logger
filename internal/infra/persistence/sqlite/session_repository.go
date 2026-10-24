// START_MODULE_CONTRACT
//   PURPOSE: Persist numbered workout sessions.
//   SCOPE: Transaction-aware session writes and deterministic workout reads.
//   DEPENDS: M-DOMAIN-SESSION
//   LINKS: M-INFRA-PERSISTENCE-SQLITE, V-M-INFRA-PERSISTENCE-SQLITE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   SessionRepository - persist independent workout sessions
//   NewSessionRepository - wire SQLite session persistence
//   Save - insert a validated session
//   NextIndex - compute maximum existing index plus one
//   FindByID - read one validated session
//   FindAll - read all sessions in descending index order
//   DeleteByID - remove a session and its sets
// END_MODULE_MAP

package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/domain/session"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

// START_BLOCK_SESSION_PERSISTENCE
type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository { return &SessionRepository{db: db} }

func (r *SessionRepository) Save(
	ctx context.Context,
	item session.Session,
) error {
	validated, err := session.NewSession(
		item.ID,
		item.Index,
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
		`INSERT INTO sessions(id, "index", created_at) VALUES (?, ?, ?)`,
		validated.ID,
		validated.Index,
		validated.CreatedAt,
	)
	return err
}

func (r *SessionRepository) NextIndex(ctx context.Context) (int, error) {
	var largest int
	err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryRowContext(
		ctx,
		`SELECT COALESCE(MAX("index"), 0) FROM sessions`,
	).Scan(&largest)
	if err != nil {
		return 0, err
	}

	if largest == int(^uint(0)>>1) {
		return 0, session.ErrInvalidIndex
	}

	return largest + 1, nil
}

func (r *SessionRepository) FindByID(
	ctx context.Context,
	id string,
) (session.Session, error) {
	var storedID string
	var index int
	var createdAt time.Time
	err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryRowContext(
		ctx,
		`SELECT id, "index", created_at FROM sessions WHERE id = ?`,
		id,
	).Scan(
		&storedID,
		&index,
		&createdAt,
	)
	if err != nil {
		return session.Session{}, err
	}

	return session.NewSession(
		storedID,
		index,
		&createdAt,
	)
}

func (r *SessionRepository) FindAll(
	ctx context.Context,
) ([]session.Session, error) {
	rows, err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).QueryContext(
		ctx,
		`SELECT id, "index", created_at FROM sessions ORDER BY "index" DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []session.Session{}
	for rows.Next() {
		var id string
		var index int
		var createdAt time.Time
		if err := rows.Scan(
			&id,
			&index,
			&createdAt,
		); err != nil {
			return nil, err
		}
		item, err := session.NewSession(
			id,
			index,
			&createdAt,
		)
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

func (r *SessionRepository) DeleteByID(
	ctx context.Context,
	id string,
) error {
	_, err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).ExecContext(
		ctx,
		`DELETE FROM sessions WHERE id = ?`,
		id,
	)
	return err
}

var _ session.SessionRepository = (*SessionRepository)(nil)

// END_BLOCK_SESSION_PERSISTENCE
