package stl

type Map[K any, V any] interface {
	Container
	// Has 检查 Map 中是否存在指定的键
	Has(K) bool
	// Find 返回 Map 中指定键对应的值的引用
	Find(K) *V
	// Insert 往 Map 中插入一个不存在的键值对或更新已存在的键对应的值
	Insert(K, V)
	// Remove 从 Map 中移除指定的键及其对应的值，返回原 Map 中是否存在指定键
	Remove(K) bool
	// ForEach 迭代 Map，对每对键值对执行指定的函数
	//  * 无法对值进行修改
	ForEach(func(K, V))
	// ForEachIf 迭代 Map，对每对键值对执行指定的函数
	//  * 当 func(K, V) bool 返回 false 时停止
	//  * 无法对值进行修改
	ForEachIf(func(K, V) bool)
	// ForEachMutable 迭代 Map，对每对键值对执行指定的函数
	//  * 可对值进行修改
	ForEachMutable(func(K, *V))
	// ForEachMutableIf 迭代 Map，对每对键值对执行指定的函数
	//  * 当 func(K, *V) bool 返回 false 时停止
	//  * 可对值进行修改
	ForEachMutableIf(func(K, *V) bool)
}

type SortedMap[K any, V any] interface {
	Map[K, V]
	// LowerBound 返回一个迭代器，指向容器中第一个满足 ele.key >= K 的元素，如果不存在，则返回 end 迭代器
	LowerBound(K) MutableMapIterator[K, V]
	// UpperBound 返回一个迭代器，指向容器中第一个满足 ele.key < K 的元素，如果不存在，则返回 end 迭代器
	UpperBound(K) MutableMapIterator[K, V]
	// FindRange 返回一个迭代器，指向 Map 中键位于 [first, last) 的区间
	FindRange(K, K) MutableMapIterator[K, V]
}
