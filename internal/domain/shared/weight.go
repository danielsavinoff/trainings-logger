package shared

import (
	"math"
)

type Weight struct {
	Value float64    `json:"value"`
	Unit  WeightUnit `json:"unit"`
}

func NewWeight(value float64, unit WeightUnit) (Weight, error) {
	if unit == "" {
		return Weight{}, ErrEmptyWeightUnit
	}

	if unit != WeightUnitG && unit != WeightUnitKg && unit != WeightUnitLb {
		return Weight{}, ErrInvalidWeightUnit
	}

	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return Weight{}, ErrInvalidWeight
	}

	return Weight{
		Value: value,
		Unit:  unit,
	}, nil
}

func (w Weight) Grams() (float64, error) {
	switch w.Unit {
	case WeightUnitG:
		return w.Value, nil

	case WeightUnitKg:
		return w.Value * 1000, nil

	case WeightUnitLb:
		return w.Value * 453.59237, nil

	default:
		return 0, ErrInvalidWeightUnit
	}
}

func (w Weight) To(unit WeightUnit) (Weight, error) {
	grams, err := w.Grams()
	if err != nil {
		return Weight{}, err
	}

	switch unit {
	case WeightUnitG:
		return Weight{Value: grams, Unit: unit}, nil

	case WeightUnitKg:
		return Weight{Value: grams / 1000, Unit: unit}, nil

	case WeightUnitLb:
		return Weight{Value: grams / 453.59237, Unit: unit}, nil

	default:
		return Weight{}, ErrInvalidWeightUnit
	}
}
