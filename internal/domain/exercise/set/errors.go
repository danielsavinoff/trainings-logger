package set

import "errors"

var (
	ErrEmptySetID            = errors.New("set id cannot be empty")
	ErrExerciseOfSetNotFound = errors.New("exercise of set not found")
	ErrSetNotFound           = errors.New("set not found")
	ErrInvalidRestInSeconds  = errors.New(
		"rest must be a nonnegative number of seconds",
	)
)
