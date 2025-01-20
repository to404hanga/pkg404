package queue

import (
	"github.com/to404hanga/pkg404/stl"
	"github.com/to404hanga/pkg404/stl/internal/heap"
)

type PriorityQueue[T any] struct {
	*heap.MinHeap[T]
}

var _ stl.PriorityQueue[any] = (*PriorityQueue[any])(nil)

func NewPriorityQueue[T stl.Ordered]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap.NewMinHeap([]T{}),
	}
}

func NewPriorityQueueFunc[T any](less stl.LessFunc[T]) *PriorityQueue[T] {
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
