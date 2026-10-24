package state

import "errors"

var (
	ErrInvalidPerformances  = errors.New("backup contains invalid performances or references")
	ErrUnsupportedVersion   = errors.New("unsupported backup version")
	ErrInvalidExercises     = errors.New("backup contains invalid or duplicate exercises")
	ErrInvalidSetReferences = errors.New("backup contains invalid set references, IDs, or timestamps")
	ErrInvalidSessions      = errors.New("backup contains invalid sessions or exercise references")
	ErrInvalidSets          = errors.New("backup contains invalid workout sets")
	ErrInvalidPreferences   = errors.New("backup contains invalid preferences")
)
