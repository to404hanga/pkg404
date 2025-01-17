package saramax

import (
	"context"
	"encoding/json"
	"pkg404/logger"
	"time"

	"github.com/IBM/sarama"
)

type BatchHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
	l  logger.Logger
}

func NewBatchHandler[T any](fn func(msgs []*sarama.ConsumerMessage, ts []T) error, l logger.Logger) sarama.ConsumerGroupHandler {
	return &BatchHandler[T]{
		fn: fn,
		l:  l,
	}
}

var _ sarama.ConsumerGroupHandler = (*BatchHandler[any])(nil)

func (bh *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (bh *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (bh *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 超时了
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				// 解决有序性，使用 channel，同一个业务发到同一个 channel
				// batch = append(batch, msg)
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					bh.l.Error("反序列化消息体失败", logger.String("topic", msg.Topic), logger.Int32("partition", msg.Partition), logger.Int64("offset", msg.Offset), logger.Error(err))
					continue
				}
				batch = append(batch, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		err := bh.fn(batch, ts)
		if err != nil {
			bh.l.Error("BatchHandler 处理消息批次失败", logger.Error(err))
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "") // 标记已消费
		}
	}
}
