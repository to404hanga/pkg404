package stl

type Queue[T any] interface {
	Container
	// Front 返回队首元素的值
	Front() T
	// Back 返回队尾元素的值
	Back() T
	// Push 往队尾添加元素
	Push(T)
	// Pop 从队首移除元素并返回
	Pop() T
}

type Deque[T any] interface {
	Container
	// Front 返回双端队列队首元素的值
	Front() T
	// Back 返回双端队列队尾元素的值
	Back() T
	// PushFront 往队首添加元素
	PushFront(T)
	// PushBack 往队尾添加元素
	PushBack(T)
	// PopFront 从队首移除元素并返回
	PopFront() T
	// PopBack 从队尾移除元素并返回
	PopBack() T
}

type PriorityQueue[T any] interface {
	Container
	// Top 返回队首元素的值
	Top() T
	// Push 往优先队列添加元素
	Push(T)
	// Pop 从队首移除元素并返回
	Pop() T
}
