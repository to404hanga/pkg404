package heap

import (
	"github.com/to404hanga/pkg404/stl/interfaces"
	"github.com/to404hanga/pkg404/stl/internal/transform"
)

type MinHeap[T any] struct {
	heap []T
	impl minHeapImpl[T]
}

var _ interfaces.Heap[any] = (*MinHeap[any])(nil)

func NewMinHeap[T interfaces.Ordered](array []T) *MinHeap[T] {
	hp := minHeapOrdered[T]{}
	hp.impl = (minHeapImpl[T])(&hp)
	hp.heap = array
	hp.impl.Make()
	return &hp.MinHeap
}

func NewMinHeapFunc[T any](array []T, less interfaces.LessFunc[T]) *MinHeap[T] {
	hp := minHeapFunc[T]{}
	hp.less = less
	hp.impl = (minHeapImpl[T])(&hp)
	hp.heap = array
	hp.impl.Make()
	return &hp.MinHeap
}

func (m *MinHeap[T]) IsMinHeap() bool {
	return m.impl.IsMinHeap()
}

func (m *MinHeap[T]) Len() int {
	return len(m.heap)
}

func (m *MinHeap[T]) Empty() bool {
	return len(m.heap) == 0
}

func (m *MinHeap[T]) Clear() {
	transform.FillZero(m.heap)
	m.heap = m.heap[:0]
}

func (m *MinHeap[T]) Push(v T) {
	m.impl.Push(v)
}

func (m *MinHeap[T]) Pop() T {
	return m.impl.Pop()
}

func (m *MinHeap[T]) Top() T {
	return m.impl.Top()
}

func (m *MinHeap[T]) Remove(idx int) T {
	return m.impl.Remove(idx)
}

type minHeapImpl[T any] interface {
	Push(v T)
	Pop() T
	Top() T
	Make()
	Remove(idx int) T
	IsMinHeap() bool
}

type minHeapOrdered[T interfaces.Ordered] struct {
	MinHeap[T]
}

func (m *minHeapOrdered[T]) IsMinHeap() bool {
	parent := 0
	for child := 1; child < len(m.heap); child++ {
		if m.heap[parent] > m.heap[child] {
			return false
		}
		if (child & 1) == 0 {
			parent++
		}
	}
	return true
}

func (m *minHeapOrdered[T]) Make() {
	n := len(m.heap)
	for i := (n >> 1) - 1; i >= 0; i-- {
		m.heapDown(i, n)
	}
}

func (m *minHeapOrdered[T]) Push(v T) {
	m.heap = append(m.heap, v)
	m.heapUp(len(m.heap) - 1)
}

func (m *minHeapOrdered[T]) Pop() T {
	n := len(m.heap) - 1
	m.heapSwap(0, n)
	m.heapDown(0, n)
	ret := m.heap[n]
	m.heap = m.heap[:n]
	return ret
}

func (m *minHeapOrdered[T]) Top() T {
	return m.heap[0]
}

func (m *minHeapOrdered[T]) Remove(idx int) T {
	h := m.heap
	n := len(h) - 1
	if n != idx {
		m.heapSwap(idx, n)
		if !m.heapDown(idx, n) {
			m.heapUp(idx)
		}
	}
	ret := h[n]
	m.heap = h[:n]
	return ret
}

func (m *minHeapOrdered[T]) heapSwap(i, j int) {
	m.heap[i], m.heap[j] = m.heap[j], m.heap[i]
}

func (m *minHeapOrdered[T]) heapUp(i int) {
	for {
		j := (i - 1) / 2
		if i == j || !(m.heap[i] < m.heap[j]) {
			break
		}
		m.heapSwap(j, i)
		i = j
	}
}

func (m *minHeapOrdered[T]) heapDown(idx, n int) bool {
	i := idx
	for {
		j := i<<1 | 1
		if j >= n || j < 0 {
			break
		}
		k := j
		l := j + 1
		if l < n && m.heap[l] < m.heap[j] {
			k = l
		}
		if !(m.heap[k] < m.heap[i]) {
			break
		}
		m.heapSwap(i, k)
		i = k
	}
	return i > idx
}

type minHeapFunc[T any] struct {
	MinHeap[T]
	less interfaces.LessFunc[T]
}

func (m *minHeapFunc[T]) Top() T {
	return m.heap[0]
}

func (m *minHeapFunc[T]) IsMinHeap() bool {
	parent := 0
	for child := 1; child < len(m.heap); child++ {
		if m.less(m.heap[child], m.heap[parent]) {
			return false
		}
		if (child & 1) == 0 {
			parent++
		}
	}
	return true
}

func (m *minHeapFunc[T]) Remove(idx int) T {
	h := m.heap
	n := len(h) - 1
	if n != idx {
		m.heapSwap(idx, n)
		if !m.heapDown(idx, n) {
			m.heapUp(idx)
		}
	}
	ret := h[n]
	m.heap = h[:n]
	return ret
}

func (m *minHeapFunc[T]) Make() {
	n := len(m.heap)
	for i := (n >> 1) - 1; i >= 0; i-- {
		m.heapDown(i, n)
	}
}

func (m *minHeapFunc[T]) Push(v T) {
	m.heap = append(m.heap, v)
	m.heapUp(len(m.heap) - 1)
}

func (m *minHeapFunc[T]) Pop() T {
	n := len(m.heap) - 1
	m.heapSwap(0, n)
	m.heapDown(0, n)
	ret := m.heap[n]
	m.heap = m.heap[:n]
	return ret
}

func (m *minHeapFunc[T]) heapSwap(i, j int) {
	m.heap[i], m.heap[j] = m.heap[j], m.heap[i]
}

func (m *minHeapFunc[T]) heapUp(i int) {
	for {
		j := (i - 1) / 2
		if i == j || !m.less(m.heap[i], m.heap[j]) {
			break
		}
		m.heapSwap(j, i)
		i = j
	}
}

func (m *minHeapFunc[T]) heapDown(idx, n int) bool {
	i := idx
	for {
		j := i<<1 | 1
		if j >= n || j < 0 {
			break
		}
		k := j
		l := j + 1
		if l < n && m.less(m.heap[l], m.heap[j]) {
			k = l
		}
		if !m.less(m.heap[k], m.heap[i]) {
			break
		}
		m.heapSwap(i, k)
		i = k
	}
	return i > idx
}
