package heap

import (
	"github.com/to404hanga/pkg404/stl/interfaces"
	"github.com/to404hanga/pkg404/stl/internal/heap"
)

func NewMinHeap[T interfaces.Ordered](array []T) *heap.MinHeap[T] {
	return heap.NewMinHeap(array)
}

func NewMinHeapFunc[T any](array []T, less interfaces.LessFunc[T]) *heap.MinHeap[T] {
	return heap.NewMinHeapFunc(array, less)
}
