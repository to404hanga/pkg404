package db

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DB interface {
	io.Closer
	Ping(ctx context.Context) error
	Find(ctx context.Context, out any, where ...any) error
	Count(ctx context.Context, model any, query string, where ...any) (count int64, err error)
	Paging(ctx context.Context, out any, offset, limit int, order, query string, where ...any) error
	First(ctx context.Context, out any, where ...any) error
	Exist(ctx context.Context, out any, where ...any) (exist bool, err error)
	Create(ctx context.Context, obj any) error
	CreateWithTx(ctx context.Context, obj any) error
	Upsert(ctx context.Context, obj any, columns []string) error
	Update(ctx context.Context, obj any, where ...any) error
	UpdateWithTx(ctx context.Context, obj any, where ...any) error
	UpdateWithResult(ctx context.Context, obj any, where ...any) (rowsAffected int64, err error)
	Updates(ctx context.Context, obj, values any, where ...any) error
	UpdateColumn(ctx context.Context, obj any, column string, value any) error
	MustUpdate(ctx context.Context, obj any, where ...any) error
	MustUpdates(ctx context.Context, obj, values any, where ...any) error
	BatchUpdateWithTx(ctx context.Context, model, values any, query string, where ...any) error
	MustBatchUpdateWithTx(ctx context.Context, model, values any, query string, where ...any) error
	Transaction(ctx context.Context, fn func(ctx context.Context, tx DB) error) error
	ExecRaw(ctx context.Context, out any, query string, where ...any) error
	Close() error
}

var (
	ErrNoRecordAffected = fmt.Errorf("db: no record/rows affected")
)

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(errors.Cause(err), gorm.ErrRecordNotFound)
}

func IsNotUpdated(err error) bool {
	return errors.Is(errors.Cause(err), ErrNoRecordAffected)
}
