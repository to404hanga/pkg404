package stl

type Heap[T any] interface {
	Container
	Push(v T)
	Pop() T
	Remove(int) T
}
