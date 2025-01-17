package connpool

import (
	"context"
	"database/sql"
	"errors"

	"github.com/LinSrHang/pkg404/logger"
	"go.uber.org/atomic"
	"gorm.io/gorm"
)

var ErrUnknownPattern = errors.New("未知的双写模式")

const (
	PatternSrcOnly  = "src_only"
	PatternDstOnly  = "dst_only"
	PatternSrcFirst = "src_first"
	PatternDstFirst = "dst_first"
)

type DoubleWritePool struct {
	src     gorm.ConnPool
	dst     gorm.ConnPool
	pattern *atomic.String
	l       logger.Logger
}

func NewDoubleWritePool(src, dst *gorm.DB, l logger.Logger) *DoubleWritePool {
	return &DoubleWritePool{
		src:     src.ConnPool,
		dst:     dst.ConnPool,
		pattern: atomic.NewString(PatternSrcOnly),
		l:       l,
	}
}

func (d *DoubleWritePool) UpdatePattern(pattern string) error {
	switch pattern {
	case PatternSrcOnly, PatternDstOnly, PatternSrcFirst, PatternDstFirst:
		d.pattern.Store(pattern)
		return nil
	default:
		return ErrUnknownPattern
	}
}

func (d *DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly:
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{
			src:     src,
			l:       d.l,
			pattern: pattern,
		}, err
	case PatternSrcFirst:
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			src.Rollback() // 回滚 src
			d.l.Error("双写目标表开启事务失败", logger.Error(err))
			return nil, err
		}
		return &DoubleWriteTx{
			src:     src,
			dst:     dst,
			l:       d.l,
			pattern: pattern,
		}, nil
	case PatternDstOnly:
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{
			dst:     dst,
			l:       d.l,
			pattern: pattern,
		}, err
	case PatternDstFirst:
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			dst.Rollback() // 回滚 dst
			d.l.Error("双写源表开启事务失败", logger.Error(err))
			return nil, err
		}
		return &DoubleWriteTx{
			src:     src,
			dst:     dst,
			l:       d.l,
			pattern: pattern,
		}, nil
	default:
		return nil, ErrUnknownPattern
	}
}

func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// 无法返回一个双写的 sql.Stmt
	panic("双写模式不支持 PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)")
}

func (d *DoubleWritePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern.Load() {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err == nil {
			_, er := d.dst.ExecContext(ctx, query, args...)
			if er != nil {
				d.l.Error("双写写入 dst 失败", logger.Error(err), logger.String("sql", query))
			}
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err == nil {
			_, er := d.src.ExecContext(ctx, query, args...)
			if er != nil {
				d.l.Error("双写写入 src 失败", logger.Error(err), logger.String("sql", query))
			}
		}
		return res, err
	default:
		return nil, ErrUnknownPattern
	}
}

func (d *DoubleWritePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern.Load() {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, ErrUnknownPattern
	}
}

func (d *DoubleWritePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern.Load() {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		panic(ErrUnknownPattern)
	}
}

type DoubleWriteTx struct {
	src     *sql.Tx
	dst     *sql.Tx
	pattern string
	l       logger.Logger
}

func (d *DoubleWriteTx) Commit() error {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Commit()
	case PatternSrcFirst:
		err := d.src.Commit()
		if err != nil {
			return err
		}
		if d.dst != nil {
			err = d.dst.Commit()
			if err != nil {
				d.l.Error("双写提交 dst 失败", logger.Error(err))
			}
		}
		return nil
	case PatternDstOnly:
		return d.dst.Commit()
	case PatternDstFirst:
		err := d.dst.Commit()
		if err != nil {
			return err
		}
		if d.src != nil {
			err = d.src.Commit()
			if err != nil {
				d.l.Error("双写提交 src 失败", logger.Error(err))
			}
		}
		return nil
	default:
		return ErrUnknownPattern
	}
}

func (d *DoubleWriteTx) Rollback() error {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Rollback()
	case PatternSrcFirst:
		err := d.src.Rollback()
		if err != nil {
			return err
		}
		if d.dst != nil {
			err = d.dst.Rollback()
			if err != nil {
				d.l.Error("双写回滚 dst 失败", logger.Error(err))
			}
		}
		return nil
	case PatternDstOnly:
		return d.dst.Rollback()
	case PatternDstFirst:
		err := d.dst.Rollback()
		if err != nil {
			return err
		}
		if d.src != nil {
			err = d.src.Rollback()
			if err != nil {
				d.l.Error("双写回滚 src 失败", logger.Error(err))
			}
		}
		return nil
	default:
		return ErrUnknownPattern
	}
}

func (d *DoubleWriteTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// 无法返回一个双写的 sql.Stmt
	panic("双写模式不支持 PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)")
}

func (d *DoubleWriteTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err == nil {
			_, er := d.dst.ExecContext(ctx, query, args...)
			if er != nil {
				d.l.Error("双写写入 dst 失败", logger.Error(err), logger.String("sql", query))
			}
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err == nil {
			_, er := d.src.ExecContext(ctx, query, args...)
			if er != nil {
				d.l.Error("双写写入 src 失败", logger.Error(err), logger.String("sql", query))
			}
		}
		return res, err
	default:
		return nil, ErrUnknownPattern
	}
}

func (d *DoubleWriteTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, ErrUnknownPattern
	}
}

func (d *DoubleWriteTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		panic(ErrUnknownPattern)
	}
}
