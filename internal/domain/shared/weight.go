package shared

import (
	"math"
)

type Weight struct {
	Value float64    `json:"value"`
	Unit  WeightUnit `json:"unit"`
}

func NewWeight(
	value float64,
	unit WeightUnit,
) (Weight, error) {
	validatedUnit, err := NewWeightUnit(string(unit))
	if err != nil {
		return Weight{}, err
	}

	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return Weight{}, ErrInvalidWeight
	}

	return Weight{
		Value: value,
		Unit:  validatedUnit,
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
	case WeightUnitG, WeightUnitKg, WeightUnitLb:
		var value float64
		switch unit {
		case WeightUnitG:
			value = grams
		case WeightUnitKg:
			value = grams / 1000
		case WeightUnitLb:
			value = grams / 453.59237
		}
		return NewWeight(
			value,
			unit,
		)
	default:
		return Weight{}, ErrInvalidWeightUnit
	}
}
