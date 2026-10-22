package identity

import "context"

type IdentityRepository interface {
	Save(
		context.Context,
		Identity,
	) error
	FindByID(
		context.Context,
		string,
	) (Identity, error)
}
