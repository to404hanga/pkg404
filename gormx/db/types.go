package db

import "context"

type DB interface {
	Ping(ctx context.Context) error
	First(ctx context.Context, out any, query string, where ...any) error
	Find(ctx context.Context, out any, query string, where ...any) error
	Insert(ctx context.Context, data any) error
}
