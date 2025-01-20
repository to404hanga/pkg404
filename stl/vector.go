package stl

type Vector[T any] interface {
	Container
	// Cap 获取容器的容量
	Cap() int
	// Reserve 增加容器的容量到指定值
	Reserve(int)
	// Shrink 移除容器中未使用到的容量
	//  执行此函数后，Vector.Cap() == Vector.Len()
	Shrink()
	// At 获取指定位置的元素的值
	//  * 若下标超出 [0, Vector.Cap()) 范围，panic
	//  * 基于切片实现的容器可以使用 [] 进行操作
	//  * 其他实现建议使用本函数
	At(int) T
	// Set 设置指定位置元素的值
	//  * 若下标超出 [0, Vector.Cap()) 范围，panic
	//  * 基于切片实现的容器可以使用 [] 进行操作
	//  * 其他实现建议使用本函数
	Set(int, T)
	// PushBack 向容器的尾部添加元素
	PushBack(T)
	// PopBack 移除并返回容器的尾部元素
	//  * 若容器为空，则 panic
	PopBack() T
	// Back 返回容器尾部的元素的值
	Back() T
	// Insert 往容器中指定位置插入若干元素
	Insert(int, ...T)
	// Append 往容器尾部追加若干元素
	Append(...T)
	// Remove 移除容器中指定位置的元素
	Remove(int)
	// RemoveRange 移除容器中指定范围的元素
	RemoveRange(int, int)
	// RemoveIf 移除容器中符合条件的元素
	RemoveIf(func(T) bool)
	// ForEach 迭代容器，对每一个元素执行指定的函数
	//  * 无法对元素进行修改
	ForEach(func(T))
	// ForEachIf 迭代容器，对每一个元素执行指定的函数
	//  * 当 func(T) bool 返回 false 时停止
	//  * 无法对元素进行修改
	ForEachIf(func(T) bool)
	// ForEachMutable 迭代容器，对每一个元素执行指定的函数
	//  * 可对元素进行修改
	ForEachMutable(func(*T))
	// ForEachMutableIf 迭代容器，对每一个元素执行指定的函数
	//  * 当 func(*T) bool 返回 false 时停止
	//  * 可对元素进行修改
	ForEachMutableIf(func(*T) bool)
	// Iterate 返回一个可修改的迭代器
	Iterate() MutableIterator[T]
	// IterateRange 返回一个可修改的迭代器，迭代 [first,last) 范围的元素
	IterateRange(int, int) MutableIterator[T]
}
