// START_MODULE_CONTRACT
//   PURPOSE: Add a catalog exercise to a workout, returning an existing performance on repeated add.
//   SCOPE: Constructor-validated workout performance ownership and transactional persistence.
//   DEPENDS: M-DOMAIN-PERFORMANCE, M-DOMAIN-SESSION, M-DOMAIN-EXERCISE
//   LINKS: M-APP-PERFORMANCES, V-M-APP-PERFORMANCES
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   AddPerformance - carry the workout and catalog exercise IDs
//   AddPerformanceResult - return the existing or created membership
//   AddPerformanceHandler - own the idempotent add transaction
//   NewAddPerformanceHandler - wire membership and parent repositories
//   Handle - validate parents and return unique membership
// END_MODULE_MAP

package command

import (
	"context"
	"database/sql"
	"errors"
	"uuid"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/performance"
)

// START_BLOCK_ADD_PERFORMANCE
type AddPerformance struct {
	SessionID  string
	ExerciseID string
}

type AddPerformanceResult struct {
	Result performance.Performance
}

type AddPerformanceHandler struct {
	performances performance.PerformanceRepository
	sessions     session.SessionRepository
	exercises    exercise.ExerciseRepository
	transactions ports.TransactionRunner[AddPerformanceResult]
}

func NewAddPerformanceHandler(performances performance.PerformanceRepository, sessions session.SessionRepository, exercises exercise.ExerciseRepository, transactions ports.TransactionRunner[AddPerformanceResult]) AddPerformanceHandler {
	return AddPerformanceHandler{
		performances: performances,
		sessions:     sessions,
		exercises:    exercises,
		transactions: transactions,
	}
}

func (h AddPerformanceHandler) Handle(ctx context.Context, command AddPerformance) (AddPerformanceResult, error) {
	return h.transactions.Run(
		ctx,
		func(ctx context.Context) (AddPerformanceResult, error) {
			item, err := performance.NewPerformance(
				uuid.NewV4().String(),
				command.SessionID,
				command.ExerciseID,
				nil,
			)
			if err != nil {
				return AddPerformanceResult{}, err
			}

			if _, err := h.sessions.FindByID(
				ctx,
				item.SessionID,
			); err != nil {
				if errors.Is(
					err,
					sql.ErrNoRows,
				) {
					return AddPerformanceResult{}, performance.ErrSessionNotFound
				}

				return AddPerformanceResult{}, err
			}
			exists, err := h.exercises.ExistsByID(
				ctx,
				item.ExerciseID,
			)
			if err != nil {
				return AddPerformanceResult{}, err
			}

			if !exists {
				return AddPerformanceResult{}, performance.ErrExerciseNotFound
			}

			if err := h.performances.Save(
				ctx,
				item,
			); err != nil {
				return AddPerformanceResult{}, err
			}
			stored, err := h.performances.FindBySessionAndExercise(
				ctx,
				item.SessionID,
				item.ExerciseID,
			)
			return AddPerformanceResult{Result: stored}, err
		},
	)
}

// END_BLOCK_ADD_PERFORMANCE
