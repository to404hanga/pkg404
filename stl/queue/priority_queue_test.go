package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue_NewPriorityQueue(t *testing.T) {
	_ = NewPriorityQueue[int]()
}

func TestPriorityQueue_Push(t *testing.T) {
	testCases := []struct {
		name    string
		data    []int
		wantRes []int
	}{
		{
			name:    "插入单个值",
			data:    []int{1},
			wantRes: []int{1},
		},
		{
			name:    "插入多个值",
			data:    []int{3, 2, 1},
			wantRes: []int{1, 2, 3},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := NewPriorityQueue[int]()
			for _, v := range tc.data {
				q.Push(v)
			}
			res := make([]int, 0)
			for !q.Empty() {
				res = append(res, q.Top())
				q.Pop()
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestPriorityQueue_PushV2(t *testing.T) {
	q := NewPriorityQueue[int]()
	q.Push(1)
	assert.Equal(t, 1, q.Top())
	q.Push(2)
	assert.Equal(t, 1, q.Top())
	q.Push(0)
	assert.Equal(t, 0, q.Top())
}
