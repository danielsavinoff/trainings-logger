// START_MODULE_CONTRACT
//   PURPOSE: Represent validated, versioned workout snapshots without local credentials.
//   SCOPE: Snapshot construction and cross-model invariants, independent of storage.
//   DEPENDS: M-DOMAIN-EXERCISE, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-SET, M-DOMAIN-USER-PREFERENCES
//   LINKS: M-DOMAIN-STATE, V-M-DOMAIN-STATE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Snapshot - carry versioned workouts and optional preferences
//   NewSnapshot - validate and reconstruct a complete snapshot
// END_MODULE_MAP

package state

import (
	"math"

	"github.com/danielsavinoff/trainings-logger/internal/domain/session"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/performance"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/preferences"
)

// START_BLOCK_SNAPSHOT
type Snapshot struct {
	Version      Version                   `json:"version"`
	Exercises    []exercise.Exercise       `json:"exercises"`
	Sessions     []session.Session         `json:"sessions"`
	Performances []performance.Performance `json:"performances"`
	Sets         []set.Set                 `json:"sets"`
	Preferences  *preferences.Preferences  `json:"preferences"`
}

func NewSnapshot(
	version Version,
	exercises []exercise.Exercise,
	sessions []session.Session,
	performances []performance.Performance,
	sets []set.Set,
	prefs *preferences.Preferences,
) (Snapshot, error) {
	validatedVersion, err := NewVersion(int(version))
	if err != nil {
		return Snapshot{}, err
	}

	result := Snapshot{Version: validatedVersion}
	ids := make(map[string]bool)
	names := make(map[string]bool)
	for _, item := range exercises {
		validated, err := exercise.NewExercise(
			item.ID,
			item.Name,
			item.BodyWeightShare,
		)
		if err != nil || ids[item.ID] || names[validated.Name] || math.IsNaN(float64(item.BodyWeightShare)) {
			return Snapshot{}, ErrInvalidExercises
		}

		ids[item.ID], names[validated.Name] = true, true
		result.Exercises = append(
			result.Exercises,
			validated,
		)
	}

	sessionIDs := make(map[string]bool)
	sessionIndices := make(map[int]bool)
	for _, item := range sessions {
		if sessionIDs[item.ID] || sessionIndices[item.Index] || item.CreatedAt.IsZero() {
			return Snapshot{}, ErrInvalidSessions
		}
		validated, err := session.NewSession(
			item.ID,
			item.Index,
			&item.CreatedAt,
		)
		if err != nil {
			return Snapshot{}, ErrInvalidSessions
		}
		sessionIDs[item.ID], sessionIndices[item.Index] = true, true
		result.Sessions = append(
			result.Sessions,
			validated,
		)
	}

	performanceIDs := make(map[string]bool)
	pairs := make(map[[2]string]bool)
	for _, item := range performances {
		pair := [2]string{item.SessionID, item.ExerciseID}
		if !sessionIDs[item.SessionID] || !ids[item.ExerciseID] || performanceIDs[item.ID] || pairs[pair] || item.CreatedAt.IsZero() {
			return Snapshot{}, ErrInvalidPerformances
		}
		validated, err := performance.NewPerformance(
			item.ID,
			item.SessionID,
			item.ExerciseID,
			&item.CreatedAt,
		)
		if err != nil {
			return Snapshot{}, ErrInvalidPerformances
		}
		performanceIDs[item.ID], pairs[pair] = true, true
		result.Performances = append(
			result.Performances,
			validated,
		)
	}

	setIDs := make(map[string]bool)
	for _, item := range sets {
		if !performanceIDs[item.PerformanceID] || setIDs[item.ID] || item.CreatedAt.IsZero() || item.CreatedAt.Year() < 1 || item.CreatedAt.Year() > 9999 {
			return Snapshot{}, ErrInvalidSetReferences
		}

		validated, err := set.NewSet(
			item.ID,
			item.PerformanceID,
			item.Repetitions,
			item.ExerciseWeight.Weight.Value,
			item.ExerciseWeight.Weight.Unit,
			item.RestInSeconds,
			&item.CreatedAt,
		)
		if err != nil {
			return Snapshot{}, ErrInvalidSets
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
			return Snapshot{}, ErrInvalidPreferences
		}

		result.Preferences = &validated
	}

	return result, nil
}

// END_BLOCK_SNAPSHOT
