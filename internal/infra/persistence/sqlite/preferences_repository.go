package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

type PreferencesRepository struct {
	db *sql.DB
}

func NewPreferencesRepository(db *sql.DB) *PreferencesRepository {
	return &PreferencesRepository{db: db}
}

func (r *PreferencesRepository) Save(
	ctx context.Context,
	p preferences.Preferences,
) error {
	validated, err := preferences.NewPreferences(p.WeightUnit)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO preferences (id, weight_unit) VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET weight_unit = excluded.weight_unit
	`

	_, err = persistence.ExecutorFromContext(
		ctx,
		r.db,
	).ExecContext(
		ctx,
		query,
		validated.WeightUnit,
	)

	return err
}

func (r *PreferencesRepository) Find(
	ctx context.Context,
) (preferences.Preferences, error) {
	const query = "SELECT weight_unit FROM preferences WHERE id = 1"

	var unit string
	err := persistence.ExecutorFromContext(
		ctx,
		r.db,
	).
		QueryRowContext(
			ctx,
			query,
		).
		Scan(&unit)

	if errors.Is(
		err,
		sql.ErrNoRows,
	) {
		return preferences.Preferences{}, preferences.ErrPreferencesNotFound
	}

	if err != nil {
		return preferences.Preferences{}, err
	}

	return preferences.NewPreferences(shared.WeightUnit(unit))
}
