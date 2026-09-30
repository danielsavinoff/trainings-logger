package shared

import "errors"

var (
	ErrEmptyWeightUnit   = errors.New("weight unit must be specified")
	ErrInvalidWeightUnit = errors.New("invalid weight unit")
	ErrInvalidWeight     = errors.New(
		"weight must be a finite nonnegative number",
	)
)
