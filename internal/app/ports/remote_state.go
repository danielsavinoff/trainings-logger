// START_MODULE_CONTRACT
//   PURPOSE: Define validated, versioned workout backups without local credentials.
//   SCOPE: Snapshot validation and remote transfer ports.
//   DEPENDS: M-DOMAIN-EXERCISE, M-DOMAIN-SET, M-DOMAIN-USER-PREFERENCES
//   LINKS: M-APP-REMOTE-STATE, V-M-APP-REMOTE-STATE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   RemoteStateSnapshot - carry versioned workouts and optional preferences
//   NewRemoteStateSnapshot - validate and reconstruct a complete backup
//   RemoteStateStore - upload and download snapshots
// END_MODULE_MAP

package ports

import (
	"context"
	"errors"
	"math"

	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
)

// START_BLOCK_REMOTE_SNAPSHOT
const RemoteStateVersion = 1

type RemoteStateSnapshot struct {
	Version     int                      `json:"version"`
	Exercises   []exercise.Exercise      `json:"exercises"`
	Sets        []set.Set                `json:"sets"`
	Preferences *preferences.Preferences `json:"preferences"`
}

func NewRemoteStateSnapshot(
	version int,
	exercises []exercise.Exercise,
	sets []set.Set,
	prefs *preferences.Preferences,
) (RemoteStateSnapshot, error) {
	if version != RemoteStateVersion {
		return RemoteStateSnapshot{}, errors.New("unsupported backup version")
	}
	result := RemoteStateSnapshot{Version: version}
	ids := make(map[string]bool)
	names := make(map[string]bool)
	for _, item := range exercises {
		validated, err := exercise.NewExercise(
			item.ID,
			item.Name,
			item.BodyWeightShare,
		)
		if err != nil || ids[item.ID] || names[item.Name] || math.IsNaN(float64(item.BodyWeightShare)) {
			return RemoteStateSnapshot{}, errors.New("backup contains invalid or duplicate exercises")
		}
		ids[item.ID], names[item.Name] = true, true
		result.Exercises = append(
			result.Exercises,
			validated,
		)
	}
	setIDs := make(map[string]bool)
	for _, item := range sets {
		if !ids[item.ExerciseID] || setIDs[item.ID] || item.CreatedAt.IsZero() || item.CreatedAt.Year() < 1 || item.CreatedAt.Year() > 9999 {
			return RemoteStateSnapshot{}, errors.New("backup contains invalid set references, IDs, or timestamps")
		}
		validated, err := set.NewSet(
			item.ID,
			item.ExerciseID,
			item.Repetitions,
			item.ExerciseWeight.Weight.Value,
			item.ExerciseWeight.Weight.Unit,
			item.RestInSeconds,
			&item.CreatedAt,
		)
		if err != nil {
			return RemoteStateSnapshot{}, errors.New("backup contains invalid workout sets")
		}
		grams, err := validated.ExerciseWeight.Weight.Grams()
		if err != nil || math.IsInf(
			grams,
			0,
		) || math.Round(grams) >= math.Exp2(63) {
			return RemoteStateSnapshot{}, errors.New("backup weight exceeds local storage limits")
		}
		setIDs[item.ID] = true
		result.Sets = append(
			result.Sets,
			validated,
		)
	}
	if prefs != nil {
		validated, err := preferences.NewPreferences(prefs.WeightUnit)
		if err != nil {
			return RemoteStateSnapshot{}, errors.New("backup contains invalid preferences")
		}
		result.Preferences = &validated
	}

	return result, nil
}

type RemoteStateStore interface {
	Store(
		context.Context,
		RemoteStateSnapshot,
	) error
	Load(
		context.Context,
	) (RemoteStateSnapshot, error)
}

// END_BLOCK_REMOTE_SNAPSHOT
