package exercise

import (
	"math"
	"strings"
)

type Exercise struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	BodyWeightShare float32 `json:"body_weight_share"`
}

func NewExercise(
	id string,
	name string,
	bodyWeightShare float32,
) (Exercise, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Exercise{}, ErrEmptyExerciseName
	}

	if math.IsNaN(float64(bodyWeightShare)) ||
		bodyWeightShare < 0 || bodyWeightShare > 1 {
		return Exercise{}, ErrInvalidBodyWeightShare
	}

	return Exercise{
		ID:              id,
		Name:            name,
		BodyWeightShare: bodyWeightShare,
	}, nil
}
