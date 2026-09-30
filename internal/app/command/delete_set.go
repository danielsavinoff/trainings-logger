package command

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/exercise/set"
)

type DeleteSet struct {
	ID string
}

type DeleteSetHandler struct {
	setRepo      set.SetRepository
	transactions ports.TransactionRunner[ports.TransactionNoResult]
}

func NewDeleteSetHandler(
	setRepo set.SetRepository,
	transactions ports.TransactionRunner[ports.TransactionNoResult],
) DeleteSetHandler {
	return DeleteSetHandler{
		setRepo:      setRepo,
		transactions: transactions,
	}
}

func (h DeleteSetHandler) Handle(
	ctx context.Context,
	command DeleteSet,
) error {
	_, err := h.transactions.Run(ctx, func(
		ctx context.Context,
	) (ports.TransactionNoResult, error) {
		return ports.TransactionNoResult{}, h.setRepo.DeleteByID(ctx, command.ID)
	})

	return err
}
