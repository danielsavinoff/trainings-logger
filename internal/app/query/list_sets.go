package query

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
)

type ListSets struct {
	ExerciseID string
}

type ListSetsResult struct {
	Result []set.Set
}

type ListSetsHandler struct {
	setRepo      set.SetRepository
	transactions ports.TransactionRunner[ListSetsResult]
}

func NewListSetsHandler(
	setRepo set.SetRepository,
	transactions ports.TransactionRunner[ListSetsResult],
) ListSetsHandler {
	return ListSetsHandler{
		setRepo:      setRepo,
		transactions: transactions,
	}
}

func (h ListSetsHandler) Handle(
	ctx context.Context,
	query ListSets,
) (ListSetsResult, error) {
	return h.transactions.Run(ctx, func(
		ctx context.Context,
	) (ListSetsResult, error) {
		sets, err := h.setRepo.FindAllByExerciseId(ctx, query.ExerciseID)
		if err != nil {
			return ListSetsResult{}, err
		}

		return ListSetsResult{Result: sets}, nil
	})
}
