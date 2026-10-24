package performance

import "errors"

var (
	ErrInvalidID        = errors.New("performance id must be a UUID")
	ErrSessionNotFound  = errors.New("session not found")
	ErrExerciseNotFound = errors.New("catalog exercise not found")
	ErrInvalidCreatedAt = errors.New("performance creation time is invalid")
)
