package saramax

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/to404hanga/pkg404/logger"
)

type Handler[T any] struct {
	l  logger.Logger
	fn func(msg *sarama.ConsumerMessage, event T) error
}

func NewHandler[T any](l logger.Logger, fn func(msg *sarama.ConsumerMessage, event T) error) sarama.ConsumerGroupHandler {
	return &Handler[T]{
		l:  l,
		fn: fn,
	}
}

var _ sarama.ConsumerGroupHandler = (*Handler[any])(nil)

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.l.Error("反序列化消息体失败", logger.String("topic", msg.Topic), logger.Int32("partition", msg.Partition), logger.Int64("offset", msg.Offset), logger.Error(err))
		}
		err = h.fn(msg, t)
		if err != nil {
			h.l.Error("处理消息失败", logger.String("topic", msg.Topic), logger.Int32("partition", msg.Partition), logger.Int64("offset", msg.Offset), logger.Error(err))
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
