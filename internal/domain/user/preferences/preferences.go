package preferences

import "github.com/danielsavinoff/trainings-logger/internal/domain/shared"

type Preferences struct {
	WeightUnit shared.WeightUnit `json:"weight_unit"`
}

func NewPreferences(
	weightUnit shared.WeightUnit,
) (Preferences, error) {
	validatedUnit, err := shared.NewWeightUnit(string(weightUnit))
	if err != nil {
		return Preferences{}, err
	}

	return Preferences{WeightUnit: validatedUnit}, nil
}
