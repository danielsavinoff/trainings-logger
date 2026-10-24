// START_MODULE_CONTRACT
//   PURPOSE: Create a numbered session for a complete workout.
//   SCOPE: Persist a constructor-validated UUID session atomically.
//   DEPENDS: M-DOMAIN-SESSION
//   LINKS: M-APP-SESSIONS, V-M-APP-SESSIONS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   CreateSession - request an automatically numbered workout
//   CreateSessionResult - return the created domain session
//   CreateSessionHandler - own transactional session creation
//   NewCreateSessionHandler - wire session creation dependencies
//   Handle - assign the next index and save a session
// END_MODULE_MAP

package command

import (
	"context"
	"uuid"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session"
)

// START_BLOCK_CREATE_SESSION
type CreateSession struct{}

type CreateSessionResult struct {
	Result session.Session
}

type CreateSessionHandler struct {
	sessions     session.SessionRepository
	transactions ports.TransactionRunner[CreateSessionResult]
}

func NewCreateSessionHandler(
	sessions session.SessionRepository,
	transactions ports.TransactionRunner[CreateSessionResult],
) CreateSessionHandler {
	return CreateSessionHandler{
		sessions:     sessions,
		transactions: transactions,
	}
}

func (h CreateSessionHandler) Handle(
	ctx context.Context,
	_ CreateSession,
) (CreateSessionResult, error) {
	return h.transactions.Run(
		ctx,
		func(ctx context.Context) (CreateSessionResult, error) {
			index, err := h.sessions.NextIndex(ctx)
			if err != nil {
				return CreateSessionResult{}, err
			}

			item, err := session.NewSession(
				uuid.NewV4().String(),
				index,
				nil,
			)
			if err != nil {
				return CreateSessionResult{}, err
			}

			if err := h.sessions.Save(
				ctx,
				item,
			); err != nil {
				return CreateSessionResult{}, err
			}

			return CreateSessionResult{Result: item}, nil
		},
	)
}

// END_BLOCK_CREATE_SESSION
