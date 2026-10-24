package set

import (
	"strings"
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
)

type Set struct {
	ID             string         `json:"id"`
	PerformanceID  string         `json:"performance_id"`
	RestInSeconds  int            `json:"rest_in_seconds"`
	Repetitions    uint8          `json:"repetitions"`
	ExerciseWeight ExerciseWeight `json:"exercise_weight"`
	CreatedAt      time.Time      `json:"created_at"`
}

func NewSet(
	id string,
	performanceID string,
	repetitions uint8,
	value float64,
	unit shared.WeightUnit,
	restInSeconds int,
	createdAt *time.Time,
) (Set, error) {
	if strings.TrimSpace(id) == "" {
		return Set{}, ErrEmptySetID
	}

	if strings.TrimSpace(performanceID) == "" {
		return Set{}, ErrPerformanceOfSetNotFound
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
		PerformanceID:  performanceID,
		RestInSeconds:  restInSeconds,
		Repetitions:    repetitions,
		ExerciseWeight: weight,
		CreatedAt:      *createdAt,
	}, nil
}
