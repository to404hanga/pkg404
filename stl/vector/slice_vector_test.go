package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/to404hanga/pkg404/stl/interfaces"
)

func TestSliceVector_Interface(t *testing.T) {
	testCases := []struct {
		name string
		fn   func()
	}{
		{
			name: "实现 interfaces.Vector 接口",
			fn: func() {
				v := NewSliceVector[int]()
				_ = interfaces.Vector[int](&v)
			},
		},
		{
			name: "实现 interfaces.Container 接口",
			fn: func() {
				v := NewSliceVector[int]()
				_ = interfaces.Container(&v)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fn() // 调用函数并检查是否 panic 或返回错误
		})
	}
}

func TestSliceVector_NewVectorCap(t *testing.T) {
	testCases := []struct {
		name    string
		cap     int
		wantLen int
		wantCap int
	}{
		{
			name:    "零容量",
			cap:     0,
			wantLen: 0,
			wantCap: 0,
		},
		{
			name:    "一般容量",
			cap:     10,
			wantLen: 0,
			wantCap: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := NewSliceVectorCap[int](tc.cap)
			assert.Equal(t, tc.wantLen, v.Len())
			assert.Equal(t, tc.wantCap, v.Cap())
		})
	}
}

func TestSliceVector_NewSliceVectorFromSlice(t *testing.T) {
	testCases := []struct {
		name    string
		slice   []int
		wantLen int
		wantCap int
		wantVec []int
	}{
		{
			name:    "单个值",
			slice:   []int{10},
			wantLen: 1,
			wantCap: 1,
			wantVec: []int{10},
		},
		{
			name:    "多个值",
			slice:   []int{10, 20, 30},
			wantLen: 3,
			wantCap: 3,
			wantVec: []int{10, 20, 30},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := NewSliceVectorFromSlice(tc.slice...)
			assert.Equal(t, tc.wantLen, v.Len())
			assert.Equal(t, tc.wantCap, v.Cap())
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], v.At(idx))
			}
		})
	}
}

func TestSliceVector_Empty(t *testing.T) {
	testCases := []struct {
		name      string
		before    func() SliceVector[int]
		wantEmpty bool
	}{
		{
			name: "容量长度皆为空",
			before: func() SliceVector[int] {
				return NewSliceVector[int]()
			},
			wantEmpty: true,
		},
		{
			name: "容量不为空，长度为空",
			before: func() SliceVector[int] {
				return NewSliceVectorCap[int](3)
			},
			wantEmpty: true,
		},
		{
			name: "容量长度皆不为空",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			wantEmpty: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			assert.Equal(t, tc.wantEmpty, v.Empty())
		})
	}
}

func TestSliceVector_Clear(t *testing.T) {
	testCases := []struct {
		name      string
		before    func() SliceVector[int]
		wantLen   int
		wantCap   int
		wantEmpty bool
	}{
		{
			name: "原容器为空",
			before: func() SliceVector[int] {
				return NewSliceVector[int]()
			},
			wantLen:   0,
			wantCap:   0,
			wantEmpty: true,
		},
		{
			name: "原容器容量不为空，长度为空",
			before: func() SliceVector[int] {
				return NewSliceVectorCap[int](3)
			},
			wantLen:   0,
			wantCap:   3,
			wantEmpty: true,
		},
		{
			name: "原容器容量长度皆不为空",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			wantLen:   0,
			wantCap:   4,
			wantEmpty: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			v.Clear()
			assert.Equal(t, tc.wantLen, v.Len())
			assert.Equal(t, tc.wantCap, v.Cap())
			assert.Equal(t, tc.wantEmpty, v.Empty())
		})
	}
}

func TestSliceVector_Reserve(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		cap     int
		wantCap int
	}{
		{
			name: "原容器为空",
			before: func() SliceVector[int] {
				return NewSliceVector[int]()
			},
			cap:     10,
			wantCap: 10,
		},
		{
			name: "原容器容量不为空，目标值大于当前值",
			before: func() SliceVector[int] {
				return NewSliceVectorCap[int](3)
			},
			cap:     4,
			wantCap: 4,
		},
		{
			name: "原容器容量不为空，目标值不大于当前值",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			cap:     3,
			wantCap: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			v.Reserve(tc.cap)
			assert.Equal(t, tc.wantCap, v.Cap())
		})
	}
}

func TestSliceVector_Shrink(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		wantCap int
	}{
		{
			name: "原容器为空",
			before: func() SliceVector[int] {
				return NewSliceVector[int]()
			},
			wantCap: 0,
		},
		{
			name: "原容器容量不为空，长度为空",
			before: func() SliceVector[int] {
				return NewSliceVectorCap[int](3)
			},
			wantCap: 0,
		},
		{
			name: "原容器容量长度不为空，容量大于长度",
			before: func() SliceVector[int] {
				v := make([]int, 1, 3)
				v[0] = 3
				return SliceVector[int](v)
			},
			wantCap: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			v.Shrink()
			assert.Equal(t, tc.wantCap, v.Cap())
		})
	}
}

func TestSliceVector_At(t *testing.T) {
	testCases := []struct {
		name     string
		before   func() SliceVector[int]
		do       func(SliceVector[int]) (int, bool)
		wantBool bool
		wantElm  int
	}{
		{
			name: "正常获取",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (int, bool) {
				return sv.At(2), false
			},
			wantBool: false,
			wantElm:  3,
		},
		{
			name: "下标超出范围",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (int, bool) {
				ret := assert.PanicsWithValue(t, "index out of range", func() {
					sv.At(3)
				})
				return 0, ret
			},
			wantBool: true,
			wantElm:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			n, b := tc.do(v)
			assert.Equal(t, tc.wantBool, b)
			assert.Equal(t, tc.wantElm, n)
		})
	}
}

func TestSliceVector_Set(t *testing.T) {
	testCases := []struct {
		name     string
		before   func() SliceVector[int]
		do       func(SliceVector[int]) (int, bool)
		wantBool bool
		wantElm  int
	}{
		{
			name: "正常设置",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (int, bool) {
				sv.Set(2, 4)
				return sv.At(2), false
			},
			wantBool: false,
			wantElm:  4,
		},
		{
			name: "下标超出范围",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (int, bool) {
				ret := assert.PanicsWithValue(t, "index out of range", func() {
					sv.Set(3, 4)
				})
				return 0, ret
			},
			wantBool: true,
			wantElm:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			n, b := tc.do(v)
			assert.Equal(t, tc.wantBool, b)
			assert.Equal(t, tc.wantElm, n)
		})
	}
}

func TestSliceVector_PushBack(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		val     int
		wantVec []int
	}{
		{
			name: "成功添加",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			val:     4,
			wantVec: []int{1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			v.PushBack(tc.val)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], v.At(idx))
			}
		})
	}
}

func TestSliceVector_PopBack(t *testing.T) {
	testCases := []struct {
		name      string
		before    func() SliceVector[int]
		do        func(SliceVector[int]) (SliceVector[int], int, bool)
		wantPanic bool
		wantEle   int
		wantVec   []int
	}{
		{
			name: "成功移除",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (SliceVector[int], int, bool) {
				tmp := sv.PopBack()
				return sv, tmp, false
			},
			wantPanic: false,
			wantEle:   3,
			wantVec:   []int{1, 2},
		},
		{
			name: "原容器为空",
			before: func() SliceVector[int] {
				return NewSliceVector[int]()
			},
			do: func(sv SliceVector[int]) (SliceVector[int], int, bool) {
				var tmp int
				ret := assert.PanicsWithValue(t, "vector is empty", func() {
					tmp = sv.PopBack()
				})
				return sv, tmp, ret
			},
			wantPanic: true,
			wantEle:   0,
			wantVec:   []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			var ele int
			var ret bool
			v, ele, ret = tc.do(v)
			assert.Equal(t, tc.wantPanic, ret)
			assert.Equal(t, tc.wantEle, ele)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], v.At(idx))
			}
		})
	}
}

func TestSliceVector_Back(t *testing.T) {
	testCases := []struct {
		name      string
		before    func() SliceVector[int]
		do        func(SliceVector[int]) (int, bool)
		wantPanic bool
		wantEle   int
	}{
		{
			name: "成功获取",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (int, bool) {
				tmp := sv.PopBack()
				return tmp, false
			},
			wantPanic: false,
			wantEle:   3,
		},
		{
			name: "原容器为空",
			before: func() SliceVector[int] {
				return NewSliceVector[int]()
			},
			do: func(sv SliceVector[int]) (int, bool) {
				var tmp int
				ret := assert.PanicsWithValue(t, "vector is empty", func() {
					tmp = sv.PopBack()
				})
				return tmp, ret
			},
			wantPanic: true,
			wantEle:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			ele, ret := tc.do(v)
			assert.Equal(t, tc.wantPanic, ret)
			assert.Equal(t, tc.wantEle, ele)
		})
	}
}

func TestSliceVector_Append(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		appends []int
		wantVec []int
	}{
		{
			name: "成功追加，原容器不为空",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			appends: []int{4, 5, 6},
			wantVec: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "成功追加，原容器为空",
			before: func() SliceVector[int] {
				return NewSliceVector[int]()
			},
			appends: []int{4, 5, 6},
			wantVec: []int{4, 5, 6},
		},
		{
			name: "成功追加，追加内容为空",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			appends: []int{},
			wantVec: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			v.Append(tc.appends...)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], v.At(idx))
			}
		})
	}
}

func TestSliceVector_Insert(t *testing.T) {
	testCases := []struct {
		name      string
		before    func() SliceVector[int]
		do        func(SliceVector[int]) (SliceVector[int], bool)
		wantPanic bool
		wantVec   []int
	}{
		{
			name: "成功插入一个",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (SliceVector[int], bool) {
				sv.Insert(1, 3)
				return sv, false
			},
			wantPanic: false,
			wantVec:   []int{1, 3, 2, 3},
		},
		{
			name: "成功插入多个",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (SliceVector[int], bool) {
				sv.Insert(1, 3, 2, 1)
				return sv, false
			},
			wantPanic: false,
			wantVec:   []int{1, 3, 2, 1, 2, 3},
		},
		{
			name: "插入错误，下标超出范围",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (SliceVector[int], bool) {
				ret := assert.PanicsWithValue(t, "index out of range", func() {
					sv.Insert(3, 4)
				})
				return sv, ret
			},
			wantPanic: true,
			wantVec:   []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			var ret bool
			v, ret = tc.do(v)
			assert.Equal(t, tc.wantPanic, ret)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], v.At(idx))
			}
		})
	}
}

func TestSliceVector_Remove(t *testing.T) {
	testCases := []struct {
		name      string
		before    func() SliceVector[int]
		do        func(SliceVector[int]) (SliceVector[int], bool)
		wantPanic bool
		wantVec   []int
	}{
		{
			name: "成功删除",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (SliceVector[int], bool) {
				sv.Remove(1)
				return sv, false
			},
			wantPanic: false,
			wantVec:   []int{1, 3},
		},
		{
			name: "删除错误，下标超出范围",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (SliceVector[int], bool) {
				ret := assert.PanicsWithValue(t, "index out of range", func() {
					sv.Remove(3)
				})
				return sv, ret
			},
			wantPanic: true,
			wantVec:   []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			var ret bool
			v, ret = tc.do(v)
			assert.Equal(t, tc.wantPanic, ret)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], v.At(idx))
			}
		})
	}
}

func TestSliceVector_RemoveRange(t *testing.T) {
	testCases := []struct {
		name      string
		before    func() SliceVector[int]
		do        func(SliceVector[int]) (SliceVector[int], bool)
		wantPanic bool
		wantVec   []int
	}{
		{
			name: "成功删除",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			do: func(sv SliceVector[int]) (SliceVector[int], bool) {
				sv.RemoveRange(1, 3)
				return sv, false
			},
			wantPanic: false,
			wantVec:   []int{1, 4},
		},
		{
			name: "删除错误，下标超出范围",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (SliceVector[int], bool) {
				ret := assert.PanicsWithValue(t, "index out of range", func() {
					sv.RemoveRange(3, 4)
				})
				return sv, ret
			},
			wantPanic: true,
			wantVec:   []int{1, 2, 3},
		},
		{
			name: "删除错误，start 不小于 end",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3)
			},
			do: func(sv SliceVector[int]) (SliceVector[int], bool) {
				ret := assert.PanicsWithValue(t, "end should be strict greater than start", func() {
					sv.RemoveRange(2, 1)
				})
				return sv, ret
			},
			wantPanic: true,
			wantVec:   []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			var ret bool
			v, ret = tc.do(v)
			assert.Equal(t, tc.wantPanic, ret)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], v.At(idx))
			}
		})
	}
}

func TestSliceVector_RemoveIf(t *testing.T) {
	testCases := []struct {
		name      string
		before    func() SliceVector[int]
		condition func(int) bool
		wantVec   []int
	}{
		{
			name: "成功删除值为偶数的元素",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			condition: func(i int) bool {
				return i%2 == 0
			},
			wantVec: []int{1, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			v.RemoveIf(tc.condition)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], v.At(idx))
			}
		})
	}
}

func TestSliceVector_ForEach(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		do      func(SliceVector[int]) []int
		wantVec []int
	}{
		{
			name: "成功迭代",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			do: func(sv SliceVector[int]) []int {
				b := make([]int, 0)
				sv.ForEach(func(val int) {
					b = append(b, val)
				})
				return b
			},
			wantVec: []int{1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			b := tc.do(v)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], b[idx])
			}
		})
	}
}

func TestSliceVector_ForEachIf(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		do      func(SliceVector[int]) []int
		wantVec []int
	}{
		{
			name: "成功迭代",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			do: func(sv SliceVector[int]) []int {
				b := make([]int, 0)
				sv.ForEachIf(func(val int) bool {
					if val == 3 {
						return false
					}
					b = append(b, val)
					return true
				})
				return b
			},
			wantVec: []int{1, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			b := tc.do(v)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], b[idx])
			}
		})
	}
}

func TestSliceVector_ForEachMutable(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		do      func(SliceVector[int]) []int
		wantVec []int
	}{
		{
			name: "成功迭代",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			do: func(sv SliceVector[int]) []int {
				sv.ForEachMutable(func(val *int) {
					*val = 1
				})
				return sv
			},
			wantVec: []int{1, 1, 1, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			b := tc.do(v)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], b[idx])
			}
		})
	}
}

func TestSliceVector_ForEachMutableIf(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		do      func(SliceVector[int]) []int
		wantVec []int
	}{
		{
			name: "成功迭代",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			do: func(sv SliceVector[int]) []int {
				sv.ForEachMutableIf(func(val *int) bool {
					if *val == 3 {
						return false
					}
					*val = 1
					return true
				})
				return sv
			},
			wantVec: []int{1, 1, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			b := tc.do(v)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], b[idx])
			}
		})
	}
}

func TestSliceVector_Iterate(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		do      func(SliceVector[int]) ([]int, []*int)
		wantVec []int
	}{
		{
			name: "成功迭代",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			do: func(sv SliceVector[int]) ([]int, []*int) {
				b := make([]int, 0)
				pb := make([]*int, 0)
				for it := sv.Iterate(); it.NotEnd(); it.MoveToNext() {
					b = append(b, it.Value())
					pb = append(pb, it.Reference())
				}
				return b, pb
			},
			wantVec: []int{1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			b, pb := tc.do(v)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], b[idx])
				assert.Equal(t, tc.wantVec[idx], *pb[idx])
			}
		})
	}
}

func TestSliceVector_IterateRange(t *testing.T) {
	testCases := []struct {
		name    string
		before  func() SliceVector[int]
		do      func(SliceVector[int]) ([]int, []*int)
		wantVec []int
	}{
		{
			name: "成功迭代",
			before: func() SliceVector[int] {
				return NewSliceVectorFromSlice(1, 2, 3, 4)
			},
			do: func(sv SliceVector[int]) ([]int, []*int) {
				b := make([]int, 0)
				pb := make([]*int, 0)
				for it := sv.IterateRange(1, 3); it.NotEnd(); it.MoveToNext() {
					b = append(b, it.Value())
					pb = append(pb, it.Reference())
				}
				return b, pb
			},
			wantVec: []int{2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.before()
			b, pb := tc.do(v)
			for idx := range tc.wantVec {
				assert.Equal(t, tc.wantVec[idx], b[idx])
				assert.Equal(t, tc.wantVec[idx], *pb[idx])
			}
		})
	}
}
