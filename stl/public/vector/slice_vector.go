package vector

import "github.com/to404hanga/pkg404/stl/internal/vector"

// NewSliceVector 新建一个基于切片的向量
func NewSliceVector[T any]() vector.SliceVector[T] {
	return vector.NewSliceVector[T]()
}

// NewSliceVectorCap 新建一个拥有初始容量的基于切片的向量
func NewSliceVectorCap[T any](n int) vector.SliceVector[T] {
	return vector.NewSliceVectorCap[T](n)
}

// NewSliceVectorFromSlice 根据初始值初始化一个基于切片的向量
func NewSliceVectorFromSlice[T any](v ...T) vector.SliceVector[T] {
	return vector.NewSliceVectorFromSlice(v...)
}
