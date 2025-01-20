package queue

import (
	"github.com/to404hanga/pkg404/stl/heap"
	"github.com/to404hanga/pkg404/stl/interfaces"
)

type PriorityQueue[T any] struct {
	*heap.MinHeap[T]
}

var _ interfaces.PriorityQueue[any] = (*PriorityQueue[any])(nil)

func NewPriorityQueue[T interfaces.Ordered]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap.NewMinHeap([]T{}),
	}
}

func NewPriorityQueueFunc[T any](less interfaces.LessFunc[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap.NewMinHeapFunc([]T{}, less),
	}
}

func (q *PriorityQueue[T]) Len() int {
	return q.MinHeap.Len()
}

func (q *PriorityQueue[T]) Empty() bool {
	return q.MinHeap.Empty()
}

func (q *PriorityQueue[T]) Clear() {
	q.MinHeap.Clear()
}

func (q *PriorityQueue[T]) Top() T {
	return q.MinHeap.Top()
}
