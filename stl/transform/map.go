package transform

// MapFromSlice 按 fn 所定义的步骤将 SRC 类型切片转换为 map[K]V 类型字典
func MapFromSlice[SRC any, K comparable, V any](src []SRC, fn func(int, SRC) (K, V)) map[K]V {
	ret := make(map[K]V, len(src))
	for idx, val := range src {
		k, v := fn(idx, val)
		ret[k] = v
	}
	return ret
}

// MapFromMap 按 fn 所定义的步骤将 SRC_K, SRC_V 类型字典转换为 DST_K, DST_V 类型字典
func MapFromMap[SRC_K comparable, SRC_V any, DST_K comparable, DST_V any](src map[SRC_K]SRC_V, fn func(SRC_K, SRC_V) (DST_K, DST_V)) map[DST_K]DST_V {
	ret := make(map[DST_K]DST_V, len(src))
	for k, v := range src {
		dk, dv := fn(k, v)
		ret[dk] = dv
	}
	return ret
}
