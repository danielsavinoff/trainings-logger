package ports

import "context"

type TransactionNoResult struct{}

type TransactionRunner[T any] interface {
	Run(
		ctx context.Context,
		fn func(context.Context) (T, error),
	) (T, error)
}
