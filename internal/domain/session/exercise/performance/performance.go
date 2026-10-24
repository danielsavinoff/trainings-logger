// START_MODULE_CONTRACT
//   PURPOSE: Represent one catalog exercise in a workout.
//   SCOPE: Validate membership identity, parents, and UTC insertion time.
//   DEPENDS: none
//   LINKS: M-DOMAIN-PERFORMANCE, V-M-DOMAIN-PERFORMANCE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Performance - retain one performance in a workout
//   NewPerformance - validate identity, parents, and creation time
// END_MODULE_MAP

package performance

import (
	"strings"
	"time"
	"uuid"
)

// START_BLOCK_PERFORMANCE
type Performance struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	ExerciseID string    `json:"exercise_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewPerformance(id string, sessionID string, exerciseID string, createdAt *time.Time) (Performance, error) {
	if _, err := uuid.Parse(id); err != nil {
		return Performance{}, ErrInvalidID
	}

	if strings.TrimSpace(sessionID) == "" {
		return Performance{}, ErrSessionNotFound
	}

	if strings.TrimSpace(exerciseID) == "" {
		return Performance{}, ErrExerciseNotFound
	}

	at := time.Now().UTC()
	if createdAt != nil {
		at = createdAt.UTC()
	}

	if at.IsZero() || at.Year() < 1 || at.Year() > 9999 {
		return Performance{}, ErrInvalidCreatedAt
	}

	return Performance{
		ID:         id,
		SessionID:  sessionID,
		ExerciseID: exerciseID,
		CreatedAt:  at,
	}, nil
}

// END_BLOCK_PERFORMANCE
