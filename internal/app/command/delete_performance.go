// START_MODULE_CONTRACT
//   PURPOSE: Remove an exercise performance and its sets from one workout.
//   SCOPE: Constructor-validated workout performance ownership and transactional persistence.
//   DEPENDS: M-DOMAIN-PERFORMANCE
//   LINKS: M-APP-PERFORMANCES, V-M-APP-PERFORMANCES
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   DeletePerformance - carry the membership ID
//   DeletePerformanceHandler - own scoped deletion
//   NewDeletePerformanceHandler - wire membership deletion
//   Handle - delete one performance and its sets
// END_MODULE_MAP

package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/performance"
)

// START_BLOCK_DELETE_PERFORMANCE
type DeletePerformance struct {
	ID string
}

type DeletePerformanceHandler struct {
	performances performance.PerformanceRepository
	transactions ports.TransactionRunner[ports.TransactionNoResult]
}

func NewDeletePerformanceHandler(performances performance.PerformanceRepository, transactions ports.TransactionRunner[ports.TransactionNoResult]) DeletePerformanceHandler {
	return DeletePerformanceHandler{
		performances: performances,
		transactions: transactions,
	}
}

func (h DeletePerformanceHandler) Handle(ctx context.Context, command DeletePerformance) error {
	_, err := h.transactions.Run(
		ctx,
		func(ctx context.Context) (ports.TransactionNoResult, error) {
			return ports.TransactionNoResult{}, h.performances.DeleteByID(
				ctx,
				command.ID,
			)
		},
	)
	return err
}

// END_BLOCK_DELETE_PERFORMANCE
