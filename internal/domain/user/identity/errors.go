package identity

import "errors"

var (
	ErrInvalidIdentity    = errors.New("invalid identity")
	ErrInvalidProvider    = errors.New("invalid identity provider")
	ErrInvalidCredentials = errors.New("invalid identity credentials")
)
