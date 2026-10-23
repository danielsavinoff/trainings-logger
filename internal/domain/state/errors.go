package state

import "errors"

var (
	ErrUnsupportedVersion   = errors.New("unsupported backup version")
	ErrInvalidExercises     = errors.New("backup contains invalid or duplicate exercises")
	ErrInvalidSetReferences = errors.New("backup contains invalid set references, IDs, or timestamps")
	ErrInvalidSets          = errors.New("backup contains invalid workout sets")
	ErrInvalidPreferences   = errors.New("backup contains invalid preferences")
)
