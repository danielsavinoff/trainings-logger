// START_MODULE_CONTRACT
//   PURPOSE: Define backup writing for validated workout snapshots.
//   SCOPE: Snapshot backup write port.
//   DEPENDS: M-DOMAIN-STATE
//   LINKS: M-APP-REMOTE-STATE, V-M-APP-REMOTE-STATE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   BackupWriter - write a snapshot backup
// END_MODULE_MAP

package ports

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/domain/state"
)

// START_BLOCK_BACKUP_WRITER
type BackupWriter interface {
	Write(
		context.Context,
		state.Snapshot,
	) error
}

// END_BLOCK_BACKUP_WRITER
