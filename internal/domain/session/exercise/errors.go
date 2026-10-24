package exercise

import "errors"

var (
	ErrExerciseInUse          = errors.New("exercise is used in a session; remove its session performances first")
	ErrEmptyExerciseID        = errors.New("exercise id cannot be empty")
	ErrInvalidBodyWeightShare = errors.New("value must be in range of 0...1")
	ErrEmptyExerciseName      = errors.New("exercise name cannot be empty")
	ErrExerciseNameConflict   = errors.New("exercise name is already used")
	ErrExerciseNotFound       = errors.New("exercise not found")
)
