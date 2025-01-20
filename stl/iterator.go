package stl

type Iterator[T any] interface {
	// NotEnd 返回当前迭代器是否指向容器末尾
	NotEnd() bool
	// MoveToNext 将迭代器指向容器的下一个元素
	MoveToNext()
	// Value 返回迭代器所指向的元素的值
	Value() T
}

type MutableIterator[T any] interface {
	Iterator[T]
	// Reference 返回迭代器所指向的元素的引用
	Reference() *T
}

type MapIterator[K any, V any] interface {
	Iterator[V]
	// Key 返回迭代器所指向的键值对的键
	Key() K
}

type MutableMapIterator[K any, V any] interface {
	MutableIterator[V]
	// Key 返回迭代器所指向的键值对的键
	Key() K
}
