// START_MODULE_CONTRACT
//   PURPOSE: List workout sessions.
//   SCOPE: Read constructor-validated sessions in descending index order.
//   DEPENDS: M-DOMAIN-SESSION
//   LINKS: M-APP-SESSIONS, V-M-APP-SESSIONS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ListSessions - carry the empty session query
//   ListSessionsResult - return constructor-validated sessions
//   ListSessionsHandler - own a workout-session query
//   NewListSessionsHandler - wire session query dependencies
//   Handle - list sessions in deterministic order
// END_MODULE_MAP

package query

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session"
)

// START_BLOCK_LIST_SESSIONS
type ListSessions struct{}

type ListSessionsResult struct {
	Result []session.Session
}

type ListSessionsHandler struct {
	sessions     session.SessionRepository
	transactions ports.TransactionRunner[ListSessionsResult]
}

func NewListSessionsHandler(
	sessions session.SessionRepository,
	transactions ports.TransactionRunner[ListSessionsResult],
) ListSessionsHandler {
	return ListSessionsHandler{
		sessions:     sessions,
		transactions: transactions,
	}
}

func (h ListSessionsHandler) Handle(
	ctx context.Context,
	query ListSessions,
) (ListSessionsResult, error) {
	return h.transactions.Run(
		ctx,
		func(ctx context.Context) (ListSessionsResult, error) {
			items, err := h.sessions.FindAll(ctx)
			return ListSessionsResult{Result: items}, err
		},
	)
}

// END_BLOCK_LIST_SESSIONS
