package db

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type db struct {
	gdb *gorm.DB
}

func New(dialector gorm.Dialector, isDebug bool, opts ...gorm.Option) (DB, error) {
	gdb, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}
	if isDebug {
		gdb = gdb.Debug()
	}
	return &db{gdb: gdb}, nil
}

func newWrapper(gdb *gorm.DB) DB {
	return &db{gdb: gdb}
}

func (db *db) Ping(ctx context.Context) error {
	sdb, err := db.gdb.DB()
	if err != nil {
		return err
	}
	return sdb.PingContext(ctx)
}

func (db *db) Find(ctx context.Context, out any, where ...any) error {
	return errors.WithStack(db.gdb.WithContext(ctx).Find(out, where...).Error)
}

func (db *db) Count(ctx context.Context, model any, query string, where ...any) (count int64, err error) {
	if err = db.gdb.WithContext(ctx).Model(model).Where(query, where...).Count(&count).Error; err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (db *db) Paging(ctx context.Context, out any, offset, limit int, order, query string, where ...any) error {
	return errors.WithStack(db.gdb.WithContext(ctx).Order(offset).Offset(offset).Limit(limit).Where(query, where...).Find(out).Error)
}

func (db *db) First(ctx context.Context, out any, where ...any) error {
	return errors.WithStack(db.gdb.WithContext(ctx).First(out, where...).Error)
}

func (db *db) Exist(ctx context.Context, out any, where ...any) (exist bool, err error) {
	res := db.gdb.WithContext(ctx).First(out, where...)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if res.Error != nil {
		return false, errors.WithStack(err)
	}
	return true, nil
}

func (db *db) Create(ctx context.Context, obj any) error {
	return errors.WithStack(db.gdb.WithContext(ctx).Create(obj).Error)
}

func (db *db) CreateWithTx(ctx context.Context, obj any) error {
	return db.gdb.Transaction(func(tx *gorm.DB) error {
		return errors.WithStack(tx.WithContext(ctx).Create(obj).Error)
	})
}

func (db *db) Upsert(ctx context.Context, obj any, columns []string) error {
	return errors.WithStack(db.gdb.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns(columns),
	}).Create(obj).Error)
}

func (db *db) Update(ctx context.Context, obj any, where ...any) error {
	_, err := db.UpdateWithResult(ctx, obj, where...)
	return err
}

func (db *db) UpdateWithTx(ctx context.Context, obj any, where ...any) error {
	return db.gdb.Transaction(func(tx *gorm.DB) error {
		return errors.WithStack(tx.WithContext(ctx).Updates(obj).Error)
	})
}

func (db *db) UpdateColumn(ctx context.Context, obj any, column string, value any) error {
	return errors.WithStack(db.gdb.Model(obj).Update(column, value).Error)
}

func (db *db) MustUpdate(ctx context.Context, obj any, where ...any) error {
	rowsAffected, err := db.UpdateWithResult(ctx, obj, where...)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		err = errors.WithStack(ErrNoRecordAffected)
	}
	return err
}

func (db *db) UpdateWithResult(ctx context.Context, obj any, where ...any) (rowsAffected int64, err error) {
	m := db.gdb.WithContext(ctx).Model(obj)
	if len(where) > 0 {
		m = m.Where(where[0], where[1:]...)
	}
	m = m.Updates(obj)
	return m.RowsAffected, errors.WithStack(err)
}

func (db *db) Updates(ctx context.Context, obj, values any, where ...any) error {
	m := db.gdb.WithContext(ctx).Model(obj)
	if len(where) > 0 {
		m = m.Where(where[0], where[1:]...)
	}
	return m.Updates(obj).Error
}

func (db *db) MustUpdates(ctx context.Context, obj, values any, where ...any) error {
	m := db.gdb.WithContext(ctx).Model(obj)
	if len(where) > 0 {
		m = m.Where(where[0], where[1:]...)
	}
	m = m.Updates(values)
	if m.Error != nil {
		return errors.WithStack(m.Error)
	}
	if m.RowsAffected == 0 {
		return errors.WithStack(ErrNoRecordAffected)
	}
	return nil
}

func (db *db) BatchUpdateWithTx(ctx context.Context, model, values any, query string, where ...any) error {
	return errors.WithStack(db.gdb.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(model).Where(query, where...).Updates(values).Error
	}))
}

func (db *db) MustBatchUpdateWithTx(ctx context.Context, model, values any, query string, where ...any) error {
	return db.gdb.Transaction(func(tx *gorm.DB) error {
		m := tx.WithContext(ctx).Model(model).Where(query, where...).Updates(values)
		if m.Error != nil {
			return errors.WithStack(m.Error)
		}
		if m.RowsAffected == 0 {
			return errors.WithStack(ErrNoRecordAffected)
		}
		return nil
	})
}

func (db *db) Transaction(ctx context.Context, fn func(ctx context.Context, tx DB) error) error {
	var err error
	paniced := true
	tx := db.gdb.Begin()
	defer func() {
		if paniced || err != nil {
			tx.Rollback()
		}
	}()
	err = fn(ctx, newWrapper(tx))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Commit().Error
	paniced = false
	return err
}

func (db *db) ExecRaw(ctx context.Context, out any, query string, where ...any) error {
	return errors.WithStack(db.gdb.WithContext(ctx).Raw(query, where...).Scan(out).Error)
}

func (db *db) Close() error {
	sdb, err := db.gdb.DB()
	if err != nil {
		return err
	}
	return sdb.Close()
}
