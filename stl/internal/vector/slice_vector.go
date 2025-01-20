package vector

import (
	"github.com/to404hanga/pkg404/stl/interfaces"
	"github.com/to404hanga/pkg404/stl/internal/transform"
)

// SliceVector 基于切片实现的向量
type SliceVector[T any] []T

var _ interfaces.Vector[any] = (*SliceVector[any])(nil)

// NewSliceVector 新建一个基于切片的向量
func NewSliceVector[T any]() SliceVector[T] {
	return (SliceVector[T])([]T{})
}

// NewSliceVectorCap 新建一个拥有初始容量的基于切片的向量
func NewSliceVectorCap[T any](n int) SliceVector[T] {
	v := make([]T, 0, n)
	return (SliceVector[T])(v)
}

// NewSliceVectorFromSlice 根据初始值初始化一个基于切片的向量
func NewSliceVectorFromSlice[T any](v ...T) SliceVector[T] {
	return (SliceVector[T])(v)
}

// Empty 实现了 interfaces.Container 接口
//
// 返回 SliceVector 是否为空
func (v SliceVector[T]) Empty() bool {
	return len(v) == 0
}

// Len 实现了 interfaces.Container 接口
//
// 返回 SliceVector 中元素的数量
func (v SliceVector[T]) Len() int {
	return len(v)
}

// Cap 实现了 interfaces.Vector 接口
//
// 返回 SliceVector 的容量
func (v SliceVector[T]) Cap() int {
	return cap(v)
}

// Clear 实现了 interfaces.Container 接口
//
// 清空 SliceVector 但容量不变
func (v *SliceVector[T]) Clear() {
	transform.FillZero(*v)
	*v = (*v)[:0]
}

// Reserve 实现了 interfaces.Vector 接口
//
// 将 SliceVector 的容量增加到指定值
//
// 若所提供的值小于当前容量，容量不变
func (v *SliceVector[T]) Reserve(capacity int) {
	if capacity > v.Cap() {
		newV := make([]T, len(*v), capacity)
		copy(newV, *v)
		*v = newV
	}
}

// Shrink 实现了 interfaces.Vector 接口
//
// 将 SliceVector 容量缩减至与长度相等
func (v *SliceVector[T]) Shrink() {
	if cap(*v) > len(*v) {
		*v = append([]T{}, *v...)
	}
}

// At 实现了 interfaces.Vector 接口
//
// 获取 SliceVector 指定位置的值，无法获取 [len(v), cap(v)) 区间的值
//
// 超出下标范围则 panic("index out of range")
func (v SliceVector[T]) At(idx int) T {
	if idx < 0 || idx >= len(v) {
		panic("index out of range")
	}
	return v[idx]
}

// Set 实现了 interfaces.Vector 接口
//
// 设置 SliceVector 指定位置的值，无法在 [len(v), cap(v)) 区间设置值
//
// 超出下标范围则 panic("index out of range")
func (v *SliceVector[T]) Set(idx int, val T) {
	if idx < 0 || idx >= len(*v) {
		panic("index out of range")
	}
	(*v)[idx] = val
}

// PushBack 实现了 interfaces.Vector 接口
//
// 在 SliceVector 尾部插入值
func (v *SliceVector[T]) PushBack(val T) {
	*v = append(*v, val)
}

// PopBack 实现了 interfaces.Vector 接口
//
// 从 SliceVector 尾部删除值
func (v *SliceVector[T]) PopBack() T {
	if v.Len() == 0 {
		panic("vector is empty")
	}
	var zero T
	val := (*v)[v.Len()-1]
	(*v)[v.Len()-1] = zero
	*v = (*v)[:v.Len()-1]
	return val
}

// Back 实现了 interfaces.Vector 接口
//
// 返回 SliceVector 尾部的值
func (v SliceVector[T]) Back() T {
	if v.Len() == 0 {
		panic("vector is empty")
	}
	return v[v.Len()-1]
}

// Append 实现了 interfaces.Vector 接口
//
// 往 SliceVector 尾部追加若干元素
func (v *SliceVector[T]) Append(vals ...T) {
	*v = append(*v, vals...)
}

// Insert 实现了 interfaces.Vector 接口
//
// 往 SliceVector 指定位置插入若干元素
func (v *SliceVector[T]) Insert(idx int, vals ...T) {
	if idx < 0 || idx >= v.Len() {
		panic("index out of range")
	}
	cpy := *v
	total := len(cpy) + len(vals)
	if total <= cap(cpy) {
		tmp := cpy[:total]
		copy(tmp[idx+len(vals):], cpy[idx:])
		copy(tmp[idx:], vals)
		*v = tmp
		return
	}
	tmp := make([]T, total)
	copy(tmp, cpy[:idx])
	copy(tmp[idx:], vals)
	copy(tmp[idx+len(vals):], cpy[idx:])
	*v = tmp
}

// Remove 实现了 interfaces.Vector 接口
//
// 删除 SliceVector 指定下标的元素
func (v *SliceVector[T]) Remove(idx int) {
	v.RemoveRange(idx, idx+1)
}

// RemoveRange 实现了 interfaces.Vector 接口
//
// 删除 SliceVector 中位于 [start, end) 区间的元素
func (v *SliceVector[T]) RemoveRange(start, end int) {
	if start < 0 || start >= len(*v) || end < 0 || end > len(*v) {
		panic("index out of range")
	}
	if start >= end {
		panic("end should be strict greater than start")
	}
	oldV := *v
	*v = append((*v)[:start], (*v)[end:]...)
	transform.FillZero(oldV[v.Len():])
}

// RemoveIf 实现了 interfaces.Vector 接口
//
// 删除 SliceVector 满足条件的元素
func (v *SliceVector[T]) RemoveIf(condition func(T) bool) {
	oldV := *v
	// INFO RemoveIf
	*v = func(slice []T, condition func(T) bool) []T {
		j := 0
		for _, v := range slice {
			if !condition(v) {
				slice[j] = v
				j++
			}
		}
		return slice[:j]
	}(*v, condition)
	transform.FillZero(oldV[v.Len():])
}

func (v SliceVector[T]) ForEach(cb func(val T)) {
	for _, val := range v {
		cb(val)
	}
}

func (v SliceVector[T]) ForEachIf(cb func(val T) bool) {
	for _, val := range v {
		if !cb(val) {
			break
		}
	}
}

func (v *SliceVector[T]) ForEachMutable(cb func(val *T)) {
	for i := range *v {
		cb(&(*v)[i])
	}
}

func (v *SliceVector[T]) ForEachMutableIf(cb func(val *T) bool) {
	for i := range *v {
		if !cb(&(*v)[i]) {
			break
		}
	}
}

func (v SliceVector[T]) Iterate() interfaces.MutableIterator[T] {
	return &vectorIterator[T]{
		vec: v,
		idx: 0,
	}
}

func (v SliceVector[T]) IterateRange(start, end int) interfaces.MutableIterator[T] {
	return &vectorIterator[T]{
		vec: v[start:end],
		idx: 0,
	}
}

type vectorIterator[T any] struct {
	vec SliceVector[T]
	idx int
}

var _ interfaces.MutableIterator[any] = (*vectorIterator[any])(nil)

func (it vectorIterator[T]) Value() T {
	return it.vec[it.idx]
}

func (it vectorIterator[T]) Reference() *T {
	return &it.vec[it.idx]
}

func (it *vectorIterator[T]) MoveToNext() {
	it.idx++
}

func (it vectorIterator[T]) NotEnd() bool {
	return it.idx < len(it.vec)
}
