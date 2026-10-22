package ports

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/domain/user/identity"
)

type LinkedIdentity struct {
	Identity    identity.Identity
	Credentials identity.IdentityCredentials
}

type IdentityAuthorization struct {
	URL           string
	BrowserOpened bool
	Complete      func(context.Context, string) (LinkedIdentity, error)
}

type IdentityLinker interface {
	Begin(
		context.Context,
		identity.Provider,
	) (IdentityAuthorization, error)
}
