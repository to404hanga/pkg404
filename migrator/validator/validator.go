package validator

import (
	"context"
	"time"

	"github.com/to404hanga/pkg404/logger"
	"github.com/to404hanga/pkg404/migrator"
	"github.com/to404hanga/pkg404/migrator/events"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Validator[T migrator.Entity] struct {
	baseValidator
	batchSize     int
	updateTime    int64
	sleepInterval time.Duration
}

func NewValidator[T migrator.Entity](base, target *gorm.DB, direction string, l logger.Logger, producer events.Producer) *Validator[T] {
	return &Validator[T]{
		baseValidator: baseValidator{
			base:      base,
			target:    target,
			direction: direction,
			l:         l,
			producer:  producer,
		},
		batchSize: 100,
		// 默认全量校验，数据没了就结束
		sleepInterval: 0,
	}
}

func (v *Validator[T]) UpdateTime(updateTime int64) *Validator[T] {
	v.updateTime = updateTime
	return v
}

func (v *Validator[T]) SleepInterval(sleepInterval time.Duration) *Validator[T] {
	v.sleepInterval = sleepInterval
	return v
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		return v.baseToTarget(ctx)
	})
	eg.Go(func() error {
		return v.targetToBase(ctx)
	})
	return eg.Wait()
}

func (v *Validator[T]) baseToTarget(ctx context.Context) error {
	offset := 0
	const limit = 100
	for {
		var srcs []T
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		err := v.base.WithContext(dbCtx).Order("id").Where("update_time >= ?", v.updateTime).Offset(offset).Limit(limit).Find(&srcs).Error
		cancel()
		switch err {
		case context.Canceled, context.DeadlineExceeded:
			return err
		case nil:
			if len(srcs) == 0 {
				return nil
			}
			err = v.dstDiff(srcs)
			if err != nil {
				return err
			}
		default:
			v.l.Error("src => dst 查询源表失败", logger.Error(err))
		}
		if len(srcs) < limit {
			return nil
		}
		offset += len(srcs)
	}
}

func (v *Validator[T]) dstDiff(srcs []T) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ids := make([]int64, 0, len(srcs))
	for _, src := range srcs {
		ids = append(ids, src.ID())
	}
	var dsts []T
	err := v.target.WithContext(ctx).Where("id IN ?", ids).Find(&dsts).Error
	if err != nil {
		return err
	}
	dstMap := v.toMap(dsts)
	for _, src := range srcs {
		dst, ok := dstMap[src.ID()]
		if !ok {
			v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
			continue
		}
		if !src.CompareTo(dst) {
			v.notify(src.ID(), events.InconsistentEventTypeNEQ)
		}
	}
	return nil
}

func (v *Validator[T]) toMap(data []T) map[int64]T {
	res := make(map[int64]T, len(data))
	for _, item := range data {
		res[item.ID()] = item
	}
	return res
}

func (v *Validator[T]) targetToBase(ctx context.Context) error {
	offset := 0
	for {
		var ts []T
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		err := v.target.WithContext(dbCtx).Order("id").Offset(offset).Limit(v.batchSize).Find(&ts).Error
		cancel()
		switch err {
		case gorm.ErrRecordNotFound:
			if v.sleepInterval > 0 {
				time.Sleep(v.sleepInterval)
				continue
			}
		case context.Canceled, context.DeadlineExceeded:
			return nil
		case nil:
			v.srcMissingRecords(ctx, ts)
		default:
			v.l.Error("dst => src 查询目标表失败", logger.Error(err))
		}
		if len(ts) < v.batchSize {
			return nil
		}
		offset += v.batchSize
	}
}

func (v *Validator[T]) srcMissingRecords(ctx context.Context, ts []T) {
	ids := make([]int64, 0, len(ts))
	for _, t := range ts {
		ids = append(ids, t.ID())
	}
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	base := v.base.WithContext(dbCtx)
	var srcTs []T
	err := base.Select("id").Where("id IN ?", ids).Find(&srcTs).Error
	switch err {
	case gorm.ErrRecordNotFound:
		v.notifySrcMissing(ts)
	case nil:
		// 计算差集
		missing := make([]T, 0)
		tmp := map[int64]struct{}{}
		for _, t := range ts {
			if _, ok := tmp[t.ID()]; !ok {
				tmp[t.ID()] = struct{}{}
			}
		}
		for _, t := range srcTs {
			if _, ok := tmp[t.ID()]; !ok {
				missing = append(missing, t)
			}
		}
		v.notifySrcMissing(missing)
	default:
		v.l.Error("dst => src 查询源表失败", logger.Error(err))
	}
}

func (v *Validator[T]) notifySrcMissing(ts []T) {
	for _, t := range ts {
		v.notify(t.ID(), events.InconsistentEventTypeBaseMissing)
	}
}
