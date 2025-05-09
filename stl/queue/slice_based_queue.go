package queue

import (
	"github.com/to404hanga/pkg404/gotools/transform"
	"github.com/to404hanga/pkg404/stl/interfaces"
)

type SliceBasedQueue[T any] struct {
	data  []T
	limit int
}

var _ interfaces.Queue[any] = (*SliceBasedQueue[any])(nil)

// NewSliceBasedQueue 新建一个基于切片实现的队列，不限制大小
func NewSliceBasedQueue[T any]() *SliceBasedQueue[T] {
	return &SliceBasedQueue[T]{
		data: make([]T, 0),
	}
}

// NewSliceBasedQueueWithLimit 新建一个基于切片实现的队列，并指定长度限制
func NewSliceBasedQueueWithLimit[T any](limit int) *SliceBasedQueue[T] {
	if limit <= 0 {
		panic("limit must be positive")
	}
	return &SliceBasedQueue[T]{
		data:  make([]T, 0, limit),
		limit: limit,
	}
}

// Empty 返回队列是否为空
func (q *SliceBasedQueue[T]) Empty() bool {
	return len(q.data) == 0
}

// Len 返回队列的长度
func (q *SliceBasedQueue[T]) Len() int {
	return len(q.data)
}

// Clear 清空队列
func (q *SliceBasedQueue[T]) Clear() {
	transform.FillZero(q.data)
	q.data = q.data[:0]
}

// Front 返回队首元素的值
func (q *SliceBasedQueue[T]) Front() T {
	if q.Empty() {
		panic("queue is empty")
	}
	return q.data[0]
}

// Push 入队
func (q *SliceBasedQueue[T]) Push(data T) {
	if q.limit > 0 && len(q.data) >= q.limit {
		panic("queue is full")
	}
	q.data = append(q.data, data)
}

// Pop 出队
func (q *SliceBasedQueue[T]) Pop() T {
	if q.Empty() {
		panic("queue is empty")
	}
	res := q.data[0]
	q.data = q.data[1:]
	return res
}
