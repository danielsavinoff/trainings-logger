package session

import "context"

type SessionRepository interface {
	Save(
		ctx context.Context,
		session Session,
	) error
	FindByID(
		ctx context.Context,
		id string,
	) (Session, error)
	NextIndex(ctx context.Context) (int, error)
	FindAll(ctx context.Context) ([]Session, error)
	DeleteByID(
		ctx context.Context,
		id string,
	) error
}
