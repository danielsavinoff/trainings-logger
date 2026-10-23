// START_MODULE_CONTRACT
//   PURPOSE: Validate and atomically import a workout snapshot.
//   SCOPE: Replace exercises, sets, and optional preferences while retaining local identities.
//   DEPENDS: M-APP-REMOTE-STATE, M-DOMAIN-STATE
//   LINKS: M-INFRA-PERSISTENCE-SQLITE, V-M-INFRA-PERSISTENCE-SQLITE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   StateImporter - own the snapshot import transaction boundary
//   NewStateImporter - wire SQLite snapshot import
//   Import - validate domain and storage constraints before atomic replacement
// END_MODULE_MAP

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/state"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
)

// START_BLOCK_STATE_IMPORT
type StateImporter struct{ db *sql.DB }

func NewStateImporter(db *sql.DB) *StateImporter {
	return &StateImporter{db: db}
}

func (r *StateImporter) Import(
	ctx context.Context,
	snapshot state.Snapshot,
) error {
	validated, err := state.NewSnapshot(
		snapshot.Version,
		snapshot.Exercises,
		snapshot.Sets,
		snapshot.Preferences,
	)
	if err != nil {
		return err
	}

	for _, item := range validated.Sets {
		grams, err := item.ExerciseWeight.Weight.Grams()
		if err != nil || math.IsInf(
			grams,
			0,
		) || math.Round(grams) >= math.Exp2(63) {
			return errors.New("backup weight exceeds local storage limits")
		}
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

var _ ports.StateImporter = (*StateImporter)(nil)

// END_BLOCK_STATE_IMPORT
