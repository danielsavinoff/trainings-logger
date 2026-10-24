package session

import "errors"

var (
	ErrInvalidSessionID = errors.New("session id must be a UUID")
	ErrInvalidIndex     = errors.New("session index must be a positive integer")
	ErrInvalidCreatedAt = errors.New("session creation time is invalid")
	ErrSessionNotFound  = errors.New("session not found")
)
