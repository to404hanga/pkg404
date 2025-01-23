package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_SliceFromSlice(t *testing.T) {
	{
		name := "int转int64"
		src := []int{1, 2, 3}
		fn := func(src int) int64 {
			return int64(src)
		}
		wantDst := []int64{1, 2, 3}

		t.Run(name, func(t *testing.T) {
			actualDst := SliceFromSlice(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
	{
		name := "int转float32"
		src := []int{1, 2, 3}
		fn := func(src int) float32 {
			return float32(src)
		}
		wantDst := []float32{1, 2, 3}

		t.Run(name, func(t *testing.T) {
			actualDst := SliceFromSlice(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
	{
		name := "结构体提取字段"
		type TestStruct struct {
			name string
			age  int
		}
		src := []TestStruct{
			{name: "name1", age: 1},
			{name: "name3", age: 3},
			{name: "name2", age: 2},
		}
		fn := func(src TestStruct) string {
			return src.name
		}
		wantDst := []string{"name1", "name3", "name2"}

		t.Run(name, func(t *testing.T) {
			actualDst := SliceFromSlice(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
	{
		name := "结构体转结构体"
		type TestStruct1 struct {
			name string
			age  int
		}
		type TestStruct2 struct {
			age2 int
		}
		src := []TestStruct1{
			{name: "name1", age: 1},
			{name: "name3", age: 3},
			{name: "name2", age: 2},
		}
		fn := func(src TestStruct1) TestStruct2 {
			var ret TestStruct2
			ret.age2 = src.age * src.age
			return ret
		}
		wantDst := []TestStruct2{
			{age2: 1},
			{age2: 9},
			{age2: 4},
		}

		t.Run(name, func(t *testing.T) {
			actualDst := SliceFromSlice(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
}

func TestSlice_SliceFromMap(t *testing.T) {
	{
		name := "提取key"
		src := map[int]int{
			1: 3,
			2: 2,
			3: 1,
		}
		fn := func(k, v int) int {
			return k
		}
		wantDst := []int{1, 2, 3}

		t.Run(name, func(t *testing.T) {
			actualDst := SliceFromMap(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
	{
		name := "提取value"
		src := map[int]int{
			1: 3,
			2: 2,
			3: 1,
		}
		fn := func(k, v int) int {
			return v
		}
		wantDst := []int{3, 2, 1}

		t.Run(name, func(t *testing.T) {
			actualDst := SliceFromMap(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
	{
		name := "kv相加"
		src := map[int]int{
			1: 3,
			2: 2,
			3: 1,
		}
		fn := func(k, v int) int {
			return k + v
		}
		wantDst := []int{4, 4, 4}

		t.Run(name, func(t *testing.T) {
			actualDst := SliceFromMap(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
}
