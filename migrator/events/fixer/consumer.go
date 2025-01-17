package fixer

import (
	"context"
	"errors"
	"time"

	"github.com/IBM/sarama"
	"github.com/LinSrHang/pkg404/logger"
	"github.com/LinSrHang/pkg404/migrator"
	"github.com/LinSrHang/pkg404/migrator/events"
	"github.com/LinSrHang/pkg404/migrator/fixer"
	"github.com/LinSrHang/pkg404/saramax"
	"gorm.io/gorm"
)

type Consumer[T migrator.Entity] struct {
	client   sarama.Client
	l        logger.Logger
	srcFirst *fixer.OverrideFixer[T]
	dstFirst *fixer.OverrideFixer[T]
	topic    string
}

func NewConsumer[T migrator.Entity](client sarama.Client, l logger.Logger, topic string, src, dst *gorm.DB) (*Consumer[T], error) {
	srcFirst, err := fixer.NewOverrideFixer[T](src, dst)
	if err != nil {
		return nil, err
	}
	dstFirst, err := fixer.NewOverrideFixer[T](dst, src)
	if err != nil {
		return nil, err
	}
	return &Consumer[T]{
		client:   client,
		l:        l,
		srcFirst: srcFirst,
		dstFirst: dstFirst,
		topic:    topic,
	}, nil
}

func (c *Consumer[T]) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("migrator-fix", c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{c.topic}, saramax.NewHandler[events.InconsistentEvent](c.l, c.Consume))
		if err != nil {
			c.l.Error("退出消费循环异常", logger.Error(err))
		}
	}()
	return nil
}

func (c *Consumer[T]) Consume(msg *sarama.ConsumerMessage, t events.InconsistentEvent) error {
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	switch t.Direction {
	case "SRC":
		return c.srcFirst.Fix(t)
	case "DST":
		return c.dstFirst.Fix(t)
	}
	return errors.New("未知的校验方向")
}
