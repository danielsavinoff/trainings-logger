// START_MODULE_CONTRACT
//   PURPOSE: Save and restore complete workout snapshots.
//   SCOPE: Local export, remote transfer, validated atomic restore.
//   DEPENDS: M-DOMAIN-EXERCISE, M-DOMAIN-SET, M-DOMAIN-USER-PREFERENCES
//   LINKS: M-APP-REMOTE-STATE, V-M-APP-REMOTE-STATE
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   RemoteStateHandler - wire local and remote snapshot ports
//   NewRemoteStateHandler - construct transfer handlers
//   Save - export consistently before uploading
//   Load - download and validate before replacing local state
// END_MODULE_MAP

package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
)

// START_BLOCK_REMOTE_ACTIONS
type RemoteStateHandler struct {
	local  ports.RemoteStateRepository
	remote ports.RemoteStateStore
}

func NewRemoteStateHandler(
	local ports.RemoteStateRepository,
	remote ports.RemoteStateStore,
) RemoteStateHandler {
	return RemoteStateHandler{local: local, remote: remote}
}

func (h RemoteStateHandler) Save(ctx context.Context) error {
	snapshot, err := h.local.Export(ctx)
	if err != nil {
		return err
	}
	return h.remote.Store(
		ctx,
		snapshot,
	)
}

func (h RemoteStateHandler) Load(ctx context.Context) error {
	snapshot, err := h.remote.Load(ctx)
	if err != nil {
		return err
	}
	validated, err := ports.NewRemoteStateSnapshot(
		snapshot.Version,
		snapshot.Exercises,
		snapshot.Sets,
		snapshot.Preferences,
	)
	if err != nil {
		return err
	}
	return h.local.Replace(
		ctx,
		validated,
	)
}

// END_BLOCK_REMOTE_ACTIONS
