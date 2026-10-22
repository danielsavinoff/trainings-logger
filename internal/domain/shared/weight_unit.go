package shared

type WeightUnit string

const (
	WeightUnitG  WeightUnit = "g"
	WeightUnitLb WeightUnit = "lb"
	WeightUnitKg WeightUnit = "kg"
)

func NewWeightUnit(value string) (WeightUnit, error) {
	unit := WeightUnit(value)
	switch unit {
	case WeightUnitG, WeightUnitKg, WeightUnitLb:
		return unit, nil
	case "":
		return "", ErrEmptyWeightUnit
	default:
		return "", ErrInvalidWeightUnit
	}
}
