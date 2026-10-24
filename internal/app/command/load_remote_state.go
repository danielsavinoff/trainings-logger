// START_MODULE_CONTRACT
//   PURPOSE: Restore a complete workout snapshot from remote storage.
//   SCOPE: Remote download, snapshot validation, and atomic local restore.
//   DEPENDS: M-DOMAIN-STATE
//   LINKS: M-APP-REMOTE-STATE, V-M-APP-REMOTE-STATE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   LoadRemoteStateHandler - wire remote download and local restore ports
//   NewLoadRemoteStateHandler - construct the load handler
//   Handle - download and validate before replacing local state
// END_MODULE_MAP

package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/state"
)

// START_BLOCK_LOAD_REMOTE_STATE
type LoadRemoteStateHandler struct {
	importer ports.StateImporter
	reader   ports.BackupReader
}

func NewLoadRemoteStateHandler(
	importer ports.StateImporter,
	reader ports.BackupReader,
) LoadRemoteStateHandler {
	return LoadRemoteStateHandler{
		importer: importer,
		reader:   reader,
	}
}

func (h LoadRemoteStateHandler) Handle(ctx context.Context) error {
	snapshot, err := h.reader.Read(ctx)
	if err != nil {
		return err
	}

	validated, err := state.NewSnapshot(
		snapshot.Version,
		snapshot.Exercises,
		snapshot.Sessions,
		snapshot.Performances,
		snapshot.Sets,
		snapshot.Preferences,
	)
	if err != nil {
		return err
	}

	return h.importer.Import(
		ctx,
		validated,
	)
}

// END_BLOCK_LOAD_REMOTE_STATE
