// START_MODULE_CONTRACT
//   PURPOSE: Export a consistent snapshot of local workout state.
//   SCOPE: Exercises, sets, and optional preferences without local identities.
//   DEPENDS: M-APP-REMOTE-STATE, M-DOMAIN-STATE
//   LINKS: M-INFRA-PERSISTENCE-SQLITE, V-M-INFRA-PERSISTENCE-SQLITE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   StateExporter - own the snapshot export transaction boundary
//   NewStateExporter - wire SQLite snapshot export
//   Export - read one consistent snapshot without creating preferences
// END_MODULE_MAP

package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/state"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

// START_BLOCK_STATE_EXPORT
type StateExporter struct{ db *sql.DB }

func NewStateExporter(db *sql.DB) *StateExporter {
	return &StateExporter{db: db}
}

func (r *StateExporter) Export(ctx context.Context) (state.Snapshot, error) {
	return persistence.NewTransactionRunner[state.Snapshot](r.db).Run(
		ctx,
		func(ctx context.Context) (state.Snapshot, error) {
			exercises, err := NewExerciseRepository(r.db).FindAll(ctx)
			if err != nil {
				return state.Snapshot{}, err
			}
			sets := []set.Set{}
			for _, item := range exercises {
				values, err := NewSetRepository(r.db).FindAllByExerciseId(
					ctx,
					item.ID,
				)
				if err != nil {
					return state.Snapshot{}, err
				}
				sets = append(
					sets,
					values...,
				)
			}

			value, err := NewPreferencesRepository(r.db).Find(ctx)
			if err != nil && !errors.Is(
				err,
				preferences.ErrPreferencesNotFound,
			) {
				return state.Snapshot{}, err
			}

			var prefs *preferences.Preferences
			if err == nil {
				prefs = &value
			}

			return state.NewSnapshot(
				state.Version1,
				exercises,
				sets,
				prefs,
			)
		},
	)
}

var _ ports.StateExporter = (*StateExporter)(nil)

// END_BLOCK_STATE_EXPORT
