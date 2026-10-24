// START_MODULE_CONTRACT
//   PURPOSE: Delete a session and its sets.
//   SCOPE: Delete by identity in a transaction; persistence cascades dependent sets.
//   DEPENDS: M-DOMAIN-SESSION
//   LINKS: M-APP-SESSIONS, V-M-APP-SESSIONS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DeleteSession - carry the session ID
//   DeleteSessionHandler - own transactional session deletion
//   NewDeleteSessionHandler - wire session deletion dependencies
//   Handle - delete a session and cascade its sets
// END_MODULE_MAP

package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session"
)

// START_BLOCK_DELETE_SESSION
type DeleteSession struct {
	ID string
}

type DeleteSessionHandler struct {
	sessions     session.SessionRepository
	transactions ports.TransactionRunner[ports.TransactionNoResult]
}

func NewDeleteSessionHandler(
	sessions session.SessionRepository,
	transactions ports.TransactionRunner[ports.TransactionNoResult],
) DeleteSessionHandler {
	return DeleteSessionHandler{
		sessions:     sessions,
		transactions: transactions,
	}
}

func (h DeleteSessionHandler) Handle(
	ctx context.Context,
	command DeleteSession,
) error {
	_, err := h.transactions.Run(
		ctx,
		func(ctx context.Context) (ports.TransactionNoResult, error) {
			return ports.TransactionNoResult{}, h.sessions.DeleteByID(
				ctx,
				command.ID,
			)
		},
	)
	return err
}

// END_BLOCK_DELETE_SESSION
