// START_MODULE_CONTRACT
//   PURPOSE: List the exercises added to a workout in insertion order.
//   SCOPE: Constructor-validated workout performance ownership and transactional persistence.
//   DEPENDS: M-DOMAIN-PERFORMANCE
//   LINKS: M-APP-PERFORMANCES, V-M-APP-PERFORMANCES
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   ListPerformances - carry a workout ID
//   ListPerformancesResult - return domain performances
//   ListPerformancesHandler - own the workout-membership query
//   NewListPerformancesHandler - wire membership reads
//   Handle - list performances in insertion order
// END_MODULE_MAP

package query

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/performance"
)

// START_BLOCK_LIST_PERFORMANCES
type ListPerformances struct {
	SessionID string
}

type ListPerformancesResult struct {
	Result []performance.Performance
}

type ListPerformancesHandler struct {
	performances performance.PerformanceRepository
	transactions ports.TransactionRunner[ListPerformancesResult]
}

func NewListPerformancesHandler(performances performance.PerformanceRepository, transactions ports.TransactionRunner[ListPerformancesResult]) ListPerformancesHandler {
	return ListPerformancesHandler{
		performances: performances,
		transactions: transactions,
	}
}

func (h ListPerformancesHandler) Handle(ctx context.Context, query ListPerformances) (ListPerformancesResult, error) {
	return h.transactions.Run(
		ctx,
		func(ctx context.Context) (ListPerformancesResult, error) {
			items, err := h.performances.FindAllBySessionID(
				ctx,
				query.SessionID,
			)
			return ListPerformancesResult{Result: items}, err
		},
	)
}

// END_BLOCK_LIST_PERFORMANCES
