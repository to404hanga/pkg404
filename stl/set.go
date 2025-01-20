package stl

type Set[K any] interface {
	Container
	// Has 检查 Set 中是否存在指定的元素
	Has(K)
	// Insert 往 Set 中插入一个不存在的元素或更新已存在的元素，返回是否插入成功
	Insert(K) bool
	// InsertN 往 Set 中插入若干个不存在的元素或更新已存在的元素，返回插入成功的元素的数量
	InsertN(...K) int
	// Remove 从 Set 中指定的元素，返回是否移除成功
	Remove(K) bool
	// RemoveN 从 Set 中移除若干个元素，返回移除成功的元素的数量
	RemoveN(...K) int
	// ForEach 迭代 Set，对每一个元素执行指定的函数
	//  * 无法对元素进行修改
	ForEach(func(K))
	// ForEachIf 迭代 Set，对每一个元素执行指定的函数
	//  * 当 func(K) bool 返回 false 时停止
	//  * 无法对元素进行修改
	ForEachIf(func(K) bool)
}

type OrderedSet[K any] interface {
	Set[K]
	// LowerBound 返回一个迭代器，指向容器中第一个满足 ele.key >= K 的元素，如果不存在，则返回 end 迭代器
	LowerBound(K) Iterator[K]
	// UpperBound 返回一个迭代器，指向容器中第一个满足 ele.key < K 的元素，如果不存在，则返回 end 迭代器
	UpperBound(K) Iterator[K]
	// FindRange 返回一个迭代器，指向 Map 中键位于 [first, last) 的区间
	FindRange(K, K) Iterator[K]
}
