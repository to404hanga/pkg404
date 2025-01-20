package queue

import (
	"github.com/to404hanga/pkg404/stl/interfaces"
	"github.com/to404hanga/pkg404/stl/internal/queue"
)

type PriorityQueue[T any] queue.PriorityQueue[T]

func NewPriorityQueue[T interfaces.Ordered]() *queue.PriorityQueue[T] {
	return queue.NewPriorityQueue[T]()
}

func NewPriorityQueueFunc[T any](less interfaces.LessFunc[T]) *queue.PriorityQueue[T] {
	return queue.NewPriorityQueueFunc[T](less)
}
