package identity

import "strings"

type Identity struct {
	ID       string
	Provider Provider
}

func NewIdentity(
	id string,
	provider Provider,
) (Identity, error) {
	if strings.TrimSpace(id) == "" {
		return Identity{}, ErrInvalidIdentity
	}

	validatedProvider, err := NewProvider(string(provider))
	if err != nil {
		return Identity{}, err
	}

	return Identity{
		ID:       id,
		Provider: validatedProvider,
	}, nil
}
