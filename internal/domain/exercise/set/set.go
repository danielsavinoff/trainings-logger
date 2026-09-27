package set

import (
	"time"

	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
)

type Set struct {
	ID             string         `json:"id"`
	Repetitions    uint8          `json:"repetitions"`
	ExerciseWeight ExerciseWeight `json:"exercise_weight"`
	CreatedAt      time.Time      `json:"created_at"`
}

func NewSet(id string, repetitions uint8, value float64, unit shared.WeightUnit, createdAt *time.Time) (Set, error) {
	weight, err := NewExerciseWeight(value, unit)
	if err != nil {
		return Set{}, err
	}

	if createdAt.IsZero() {
		now := time.Now().UTC()
		createdAt = &now
	}

	return Set{
		ID:             id,
		Repetitions:    repetitions,
		ExerciseWeight: weight,
		CreatedAt:      *createdAt,
	}, nil
}
