package transform

// FillZero 使用零值填充切片
func FillZero[T any](slice []T) {
	clear(slice)
}
