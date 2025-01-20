package interfaces

type Container interface {
	// Empty 返回容器是否为空（即没有元素）
	//  等效于 Container.Len() == 0
	Empty() bool
	// Len 返回容器中元素的数量
	Len() int
	// Clear 清空容器中的所有元素
	//  调用此函数后，Container.Len() 将会返回 0
	Clear()
}
