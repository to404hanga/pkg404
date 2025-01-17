package validator

import (
	"context"

	"github.com/LinSrHang/pkg404/logger"
	"github.com/LinSrHang/pkg404/migrator"
	"github.com/LinSrHang/pkg404/migrator/events"
	"gorm.io/gorm"
)

type CanalIncrValidator[T migrator.Entity] struct {
	baseValidator
}

func NewCanalIncrValidator[T migrator.Entity](base, target *gorm.DB, direction string, l logger.Logger, producer events.Producer) *CanalIncrValidator[T] {
	return &CanalIncrValidator[T]{
		baseValidator: baseValidator{
			base:      base,
			target:    target,
			direction: direction,
			l:         l,
			producer:  producer,
		},
	}
}

func (v *CanalIncrValidator[T]) Validate(ctx context.Context, id int64) error {
	var base T
	err := v.base.WithContext(ctx).Where("id = ?", id).First(&base).Error
	switch err {
	case nil:
		var target T
		err = v.target.WithContext(ctx).Where("id =?", id).First(&target).Error
		switch err {
		case nil:
			if !base.CompareTo(target) {
				v.notify(id, events.InconsistentEventTypeNEQ)
			}
			return nil
		case gorm.ErrRecordNotFound:
			v.notify(id, events.InconsistentEventTypeTargetMissing)
			return nil
		default:
			return err
		}
	case gorm.ErrRecordNotFound:
		var target T
		err = v.target.WithContext(ctx).Where("id =?", id).First(&target).Error
		switch err {
		case nil:
			v.notify(id, events.InconsistentEventTypeBaseMissing)
			return nil
		case gorm.ErrRecordNotFound:
			return nil
		default:
			return err
		}
	default:
		// 未知错误
		return err
	}
}
