// START_MODULE_CONTRACT
//   PURPOSE: Export and atomically restore local workout backups.
//   SCOPE: Exercises, sets, and optional preferences; local identities are retained.
//   DEPENDS: M-APP-REMOTE-STATE
//   LINKS: M-INFRA-PERSISTENCE-SQLITE, V-M-INFRA-PERSISTENCE-SQLITE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   RemoteStateRepository - own snapshot persistence and transaction boundary
//   NewRemoteStateRepository - wire SQLite repositories
//   Export - read one consistent snapshot without creating preferences
//   Replace - validate then restore workout tables atomically
// END_MODULE_MAP

package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

// START_BLOCK_LOCAL_SNAPSHOTS
type RemoteStateRepository struct{ db *sql.DB }

func NewRemoteStateRepository(db *sql.DB) *RemoteStateRepository {
	return &RemoteStateRepository{db: db}
}

func (r *RemoteStateRepository) Export(ctx context.Context) (ports.RemoteStateSnapshot, error) {
	return persistence.NewTransactionRunner[ports.RemoteStateSnapshot](r.db).Run(
		ctx,
		func(ctx context.Context) (ports.RemoteStateSnapshot, error) {
			exercises, err := NewExerciseRepository(r.db).FindAll(ctx)
			if err != nil {
				return ports.RemoteStateSnapshot{}, err
			}
			sets := []set.Set{}
			for _, item := range exercises {
				values, err := NewSetRepository(r.db).FindAllByExerciseId(
					ctx,
					item.ID,
				)
				if err != nil {
					return ports.RemoteStateSnapshot{}, err
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
				return ports.RemoteStateSnapshot{}, err
			}

			var prefs *preferences.Preferences
			if err == nil {
				prefs = &value
			}

			return ports.NewRemoteStateSnapshot(
				ports.RemoteStateVersion,
				exercises,
				sets,
				prefs,
			)
		},
	)
}

func (r *RemoteStateRepository) Replace(
	ctx context.Context,
	snapshot ports.RemoteStateSnapshot,
) error {
	validated, err := ports.NewRemoteStateSnapshot(
		snapshot.Version,
		snapshot.Exercises,
		snapshot.Sets,
		snapshot.Preferences,
	)
	if err != nil {
		return err
	}
	_, err = persistence.NewTransactionRunner[ports.TransactionNoResult](r.db).Run(
		ctx,
		func(ctx context.Context) (ports.TransactionNoResult, error) {
			for _, statement := range []string{"DELETE FROM sets", "DELETE FROM exercises", "DELETE FROM preferences"} {
				if _, err := persistence.ExecutorFromContext(
					ctx,
					r.db,
				).ExecContext(
					ctx,
					statement,
				); err != nil {
					return ports.TransactionNoResult{}, err
				}
			}
			for _, item := range validated.Exercises {
				if err := NewExerciseRepository(r.db).Save(
					ctx,
					item,
				); err != nil {
					return ports.TransactionNoResult{}, err
				}
			}
			for _, item := range validated.Sets {
				if err := NewSetRepository(r.db).Save(
					ctx,
					item,
				); err != nil {
					return ports.TransactionNoResult{}, err
				}
			}
			if validated.Preferences != nil {
				if err := NewPreferencesRepository(r.db).Save(
					ctx,
					*validated.Preferences,
				); err != nil {
					return ports.TransactionNoResult{}, err
				}
			}
			return ports.TransactionNoResult{}, nil
		},
	)
	return err
}

var _ ports.RemoteStateRepository = (*RemoteStateRepository)(nil)

// END_BLOCK_LOCAL_SNAPSHOTS
