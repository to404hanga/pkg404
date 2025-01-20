package interfaces

// Signed 有符号整数类型接口
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned 无符号整数类型接口
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer 整数类型接口（Signed | Unsigned）
//  Signed（~int | ~int8 | ~int16 | ~int32 | ~int64）
//  Unsigned（~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr）
type Integer interface {
	Signed | Unsigned
}

// Float 浮点数类型接口
type Float interface {
	~float32 | ~float64
}

// Ordered 可排序类型接口（Integer | Float | ~string）
//  Integer（Signed | Unsigned）
//    Signed（~int | ~int8 | ~int16 | ~int32 | ~int64）
//    Unsigned（~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr）
//  Float（~float32 | ~float64）
type Ordered interface {
	Integer | Float | ~string
}

// Number 数字类型接口（Integer | Float）
//  Integer（Signed | Unsigned）
//    Signed（~int | ~int8 | ~int16 | ~int32 | ~int64）
//    Unsigned（~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr）
//  Float（~float32 | ~float64）
type Number interface {
	Integer | Float
}

// LessFunc 定义一个可以对 T 类型元素进行大小比较的函数，返回 left 是否小于 right
type LessFunc[T any] func(left, right T) bool

// CompareFunc 定义一个可以对 T 类型元素进行大小比较的函数，返回值列表如下
//  | 返回值 |     含义      |
//  | ----- | ------------- |
//  |   1   | left >  right |
//  |   0   | left == right |
//  |  -1   | left <  right |
type CompareFunc[T any] func(left, right T) int

// HashFunc 定义一个可以对 T 类型元素生成 uint64 类型哈希值的函数
type HashFunc[T any] func(t T) uint64
