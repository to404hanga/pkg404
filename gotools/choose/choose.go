package choose

import (
	"github.com/to404hanga/pkg404/stl/interfaces"
)

func IF[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

func MAX[T interfaces.Ordered](vals ...T) T {
	maxVal := vals[0]
	for _, val := range vals {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

func MIN[T interfaces.Ordered](vals ...T) T {
	minVal := vals[0]
	for _, val := range vals {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

func MAXBY[T any](compareFunc func(T, T) bool, vals ...T) T {
	maxVal := vals[0]
	for _, val := range vals {
		if compareFunc(val, maxVal) {
			maxVal = val
		}
	}
	return maxVal
}

func MINBY[T any](compareFunc func(T, T) bool, vals ...T) T {
	minVal := vals[0]
	for _, val := range vals {
		if compareFunc(val, minVal) {
			minVal = val
		}
	}
	return minVal
}
