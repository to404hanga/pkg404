package saramax

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/to404hanga/pkg404/logger"
)

type HandlerV1[T any] struct {
	l      logger.Logger
	fn     func(msg *sarama.ConsumerMessage, event T) error
	vector *prometheus.SummaryVec
}

func NewHandlerV1[T any](consumer string, l logger.Logger, fn func(msg *sarama.ConsumerMessage, event T) error) sarama.ConsumerGroupHandler {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "saramax",
		Subsystem: "consumer_handler",
		Name:      consumer,
	}, []string{
		"topic",
		"error",
	})
	return &HandlerV1[T]{
		l:      l,
		fn:     fn,
		vector: vector,
	}
}

var _ sarama.ConsumerGroupHandler = (*HandlerV1[any])(nil)

func (h *HandlerV1[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *HandlerV1[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}
func (h *HandlerV1[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		h.consumerClaim(msg)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (h *HandlerV1[T]) consumerClaim(msg *sarama.ConsumerMessage) {
	start := time.Now()
	var err error
	defer func() {
		errInfo := strconv.FormatBool(err != nil)
		duration := time.Since(start).Milliseconds()
		h.vector.WithLabelValues(msg.Topic, errInfo).Observe(float64(duration))
	}()
	var t T
	err = json.Unmarshal(msg.Value, &t)
	if err != nil {
		h.l.Error("反序列化消息体失败", logger.String("topic", msg.Topic), logger.Int32("partition", msg.Partition), logger.Int64("offset", msg.Offset), logger.Error(err))
	}
	err = h.fn(msg, t)
	if err != nil {
		h.l.Error("处理消息失败", logger.String("topic", msg.Topic), logger.Int32("partition", msg.Partition), logger.Int64("offset", msg.Offset), logger.Error(err))
	}
}
