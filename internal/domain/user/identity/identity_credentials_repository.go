package identity

import "context"

type IdentityCredentialsRepository interface {
	Save(
		context.Context,
		IdentityCredentials,
	) error
	FindByID(
		context.Context,
		string,
	) (IdentityCredentials, error)
	FindLatest(context.Context) (IdentityCredentials, error)
}
