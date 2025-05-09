package transform

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_MapFromSlice(t *testing.T) {
	{
		name := "int转map[int]string"
		src := []int{1, 2, 3}
		fn := func(idx int, val int) (int, string) {
			return val, strconv.Itoa(val)
		}
		wantDst := map[int]string{
			1: "1",
			2: "2",
			3: "3",
		}

		t.Run(name, func(t *testing.T) {
			actualDst := MapFromSlice(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
	{
		name := "int转map[string]float32"
		src := []int{1, 2, 3}
		fn := func(idx int, val int) (string, float32) {
			return strconv.Itoa(val), float32(val)
		}
		wantDst := map[string]float32{
			"1": 1.0,
			"2": 2.0,
			"3": 3.0,
		}

		t.Run(name, func(t *testing.T) {
			actualDst := MapFromSlice(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
	{
		name := "结构体提取字段作为key"
		type TestStruct struct {
			id   string
			name string
			age  int
		}
		src := []TestStruct{
			{id: "1", name: "name1", age: 1},
			{id: "2", name: "name3", age: 3},
			{id: "3", name: "name2", age: 2},
		}
		fn := func(idx int, val TestStruct) (string, TestStruct) {
			return val.id, val
		}
		wantDst := map[string]TestStruct{
			"1": {id: "1", name: "name1", age: 1},
			"2": {id: "2", name: "name3", age: 3},
			"3": {id: "3", name: "name2", age: 2},
		}

		t.Run(name, func(t *testing.T) {
			actualDst := MapFromSlice(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
}

func TestMap_MapFromMap(t *testing.T) {
	{
		name := "kv翻转"
		src := map[int]string{
			1: "3",
			2: "2",
			3: "1",
		}
		fn := func(k int, v string) (string, int) {
			return v, k
		}
		wantDst := map[string]int{
			"3": 1,
			"2": 2,
			"1": 3,
		}

		t.Run(name, func(t *testing.T) {
			actualDst := MapFromMap(src, fn)
			assert.Equal(t, wantDst, actualDst)
		})
	}
}
