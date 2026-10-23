package ports

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/domain/state"
)

type StateExporter interface {
	Export(context.Context) (state.Snapshot, error)
}
