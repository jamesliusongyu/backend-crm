package database

import "context"

type Collection[S any, T any] interface {
	GetAll(ctx context.Context, tenant string) ([]T, error)
	GetByID(ctx context.Context, id string, tenant string) (T, error)
	GetByKeyValue(ctx context.Context, key string, value string, tenant string) (T, error)
	GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]T, error)
	Create(ctx context.Context, entity S) (string, error)
	Update(ctx context.Context, id string, entity S) error
	Delete(ctx context.Context, id string) error
}
