package set

import (
	"strings"
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
)

type Set struct {
	ID             string         `json:"id"`
	ExerciseID     string         `json:"exercise_id"`
	RestInSeconds  int            `json:"rest_in_seconds"`
	Repetitions    uint8          `json:"repetitions"`
	ExerciseWeight ExerciseWeight `json:"exercise_weight"`
	CreatedAt      time.Time      `json:"created_at"`
}

func NewSet(
	id string,
	exerciseID string,
	repetitions uint8,
	value float64,
	unit shared.WeightUnit,
	restInSeconds int,
	createdAt *time.Time,
) (Set, error) {
	if strings.TrimSpace(id) == "" {
		return Set{}, ErrEmptySetID
	}

	if strings.TrimSpace(exerciseID) == "" {
		return Set{}, ErrExerciseOfSetNotFound
	}

	if restInSeconds < 0 {
		return Set{}, ErrInvalidRestInSeconds
	}

	weight, err := NewExerciseWeight(
		value,
		unit,
	)
	if err != nil {
		return Set{}, err
	}

	if createdAt == nil || createdAt.IsZero() {
		now := time.Now().UTC()
		createdAt = &now
	}

	return Set{
		ID:             id,
		ExerciseID:     exerciseID,
		RestInSeconds:  restInSeconds,
		Repetitions:    repetitions,
		ExerciseWeight: weight,
		CreatedAt:      *createdAt,
	}, nil
}
