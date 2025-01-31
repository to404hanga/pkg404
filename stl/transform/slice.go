package transform

// SliceFromSlice 将 SRC 类型切片按 fn 所定义的步骤转换为 DST 类型切片
func SliceFromSlice[SRC any, DST any](src []SRC, fn func(int, SRC) DST) []DST {
	ret := make([]DST, 0, len(src))
	for idx, val := range src {
		ret = append(ret, fn(idx, val))
	}
	return ret
}

// SliceFromMap 将 map[K]V 类型字典按 fn 所定义的步骤转换为 DST 类型切片
func SliceFromMap[K comparable, V any, DST any](src map[K]V, fn func(K, V) DST) []DST {
	ret := make([]DST, 0, len(src))
	for k, v := range src {
		ret = append(ret, fn(k, v))
	}
	return ret
}
