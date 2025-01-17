package validator

import (
	"context"
	"time"

	"github.com/LinSrHang/pkg404/logger"
	"github.com/LinSrHang/pkg404/migrator/events"
	"gorm.io/gorm"
)

type baseValidator struct {
	base      *gorm.DB
	target    *gorm.DB
	direction string // 告知以 src 为准还是以 dst 为准
	l         logger.Logger
	producer  events.Producer
}

// 上报不一致数据
func (v *baseValidator) notify(id int64, typ string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	evt := events.InconsistentEvent{
		Direction: v.direction,
		ID:        id,
		Type:      typ,
	}
	err := v.producer.ProduceInconsistentEvent(ctx, evt)
	if err != nil {
		v.l.Error("发送消息失败", logger.Error(err), logger.Any("event", evt))
	}
}
