package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/to404hanga/pkg404/stl"
	"github.com/to404hanga/pkg404/stl/internal"
)

func TestMinHeap_NewMinHeap(t *testing.T) {
	testCases := []struct {
		name     string
		data     []int
		wantBool bool
	}{
		{
			name:     "成功构建",
			data:     []int{5, 4, 3, 2, 1},
			wantBool: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMinHeap(tc.data)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
		})
	}
}

func TestMinHeap_NewMinHeapFunc(t *testing.T) {
	testCases := []struct {
		name     string
		data     []int
		less     stl.LessFunc[int]
		wantBool bool
	}{
		{
			name:     "成功构建",
			data:     []int{1, 2, 3, 4, 5},
			less:     internal.OrderedGreater[int],
			wantBool: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMinHeapFunc(tc.data, tc.less)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
		})
	}
}

func TestMinHeap_Push(t *testing.T) {
	testCases := []struct {
		name     string
		val      int
		wantBool bool
	}{
		{
			name:     "成功插入",
			val:      1,
			wantBool: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMinHeap([]int{})
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
			heap.Push(tc.val)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
		})
	}
}

func TestMinHeap_Pop(t *testing.T) {
	testCases := []struct {
		name     string
		data     []int
		wantEle  int
		wantBool bool
	}{
		{
			name:     "成功弹出",
			data:     []int{1},
			wantEle:  1,
			wantBool: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMinHeap(tc.data)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
			assert.Equal(t, tc.wantEle, heap.Pop())
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
		})
	}
}

func TestMinHeap_PushFunc(t *testing.T) {
	testCases := []struct {
		name     string
		val      int
		wantBool bool
	}{
		{
			name:     "成功插入",
			val:      1,
			wantBool: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMinHeapFunc([]int{}, internal.OrderedGreater)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
			heap.Push(tc.val)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
		})
	}
}

func TestMinHeap_PopFunc(t *testing.T) {
	testCases := []struct {
		name     string
		data     []int
		wantEle  int
		wantBool bool
	}{
		{
			name:     "成功弹出",
			data:     []int{1},
			wantEle:  1,
			wantBool: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMinHeapFunc(tc.data, internal.OrderedGreater)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
			assert.Equal(t, tc.wantEle, heap.Pop())
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
		})
	}
}

func TestMinHeap_Remove(t *testing.T) {
	testCases := []struct {
		name      string
		data      []int
		removeIdx int
		wantEle   int
		wantBool  bool
	}{
		{
			name:      "成功删除",
			data:      []int{5, 4, 3, 2, 1},
			removeIdx: 3,
			wantEle:   1,
			wantBool:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMinHeap(tc.data)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
			heap.Remove(tc.removeIdx)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
		})
	}
}

func TestMinHeap_RemoveFunc(t *testing.T) {
	testCases := []struct {
		name      string
		data      []int
		removeIdx int
		wantEle   int
		wantBool  bool
	}{
		{
			name:      "成功删除",
			data:      []int{5, 4, 3, 2, 1},
			removeIdx: 3,
			wantEle:   1,
			wantBool:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMinHeapFunc(tc.data, internal.OrderedGreater)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
			heap.Remove(tc.removeIdx)
			assert.Equal(t, tc.wantBool, heap.IsMinHeap())
		})
	}
}
