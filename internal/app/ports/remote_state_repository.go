package ports

import "context"

type RemoteStateRepository interface {
	Export(context.Context) (RemoteStateSnapshot, error)
	Replace(
		context.Context,
		RemoteStateSnapshot,
	) error
}
