// START_MODULE_CONTRACT
//   PURPOSE: Save a complete workout snapshot to remote storage.
//   SCOPE: Consistent local export and remote upload.
//   DEPENDS: M-DOMAIN-STATE
//   LINKS: M-APP-REMOTE-STATE, V-M-APP-REMOTE-STATE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   SaveRemoteStateHandler - wire local export and remote storage ports
//   NewSaveRemoteStateHandler - construct the save handler
//   Handle - export consistently before uploading
// END_MODULE_MAP

package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
)

// START_BLOCK_SAVE_REMOTE_STATE
type SaveRemoteStateHandler struct {
	exporter ports.StateExporter
	writer   ports.BackupWriter
}

func NewSaveRemoteStateHandler(
	exporter ports.StateExporter,
	writer ports.BackupWriter,
) SaveRemoteStateHandler {
	return SaveRemoteStateHandler{
		exporter: exporter,
		writer:   writer,
	}
}

func (h SaveRemoteStateHandler) Handle(ctx context.Context) error {
	snapshot, err := h.exporter.Export(ctx)
	if err != nil {
		return err
	}

	return h.writer.Write(
		ctx,
		snapshot,
	)
}

// END_BLOCK_SAVE_REMOTE_STATE
