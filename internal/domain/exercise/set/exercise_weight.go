package set

import "github.com/danielsavinoff/trainings-logger/internal/domain/shared"

type ExerciseWeight struct {
	Weight shared.Weight `json:"weight"`
}

func NewExerciseWeight(value float64, unit shared.WeightUnit) (ExerciseWeight, error) {
	weight, err := shared.NewWeight(value, unit)
	if err != nil {
		return ExerciseWeight{}, err
	}

	return ExerciseWeight{
		Weight: weight,
	}, nil
}
