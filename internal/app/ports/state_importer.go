package ports

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/domain/state"
)

type StateImporter interface {
	Import(
		context.Context,
		state.Snapshot,
	) error
}
