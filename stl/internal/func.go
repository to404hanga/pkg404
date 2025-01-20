package internal

import "github.com/to404hanga/pkg404/stl"

// OrderedLess 返回 stl.Ordered 类型 left < right 的真值
func OrderedLess[T stl.Ordered](left, right T) bool {
	return left < right
}

// OrderedGreater 返回 stl.Ordered 类型 left > right 的真值
func OrderedGreater[T stl.Ordered](left, right T) bool {
	return left > right
}

// OrderedCompare 提供 stl.Ordered 类型的默认 CompareFunc 函数
func OrderedCompare[T stl.Ordered](left, right T) int {
	if left == right {
		return 0
	}
	if left < right {
		return -1
	}
	return 1
}
