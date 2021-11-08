package model_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

type testEntry struct {
	val float64
}

func getVal(e interface{}) float64 {
	return e.(*testEntry).val
}

func (e *testEntry) ToString() string {
	return fmt.Sprintf("%v", e.val)
}

func TestHeap(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		entries []*testEntry
	}{
		{entries: []*testEntry{{1}, {2}, {3}, {4}, {5}}},
		{entries: []*testEntry{{5}, {4}, {3}, {2}, {1}}},
		{entries: []*testEntry{{3}, {2}, {4}, {1}, {5}}},
		{entries: []*testEntry{{1}, {3}, {5}, {4}, {2}}},
	}

	for i, tc := range testCases {
		h := model.NewHeap(getVal)
		assert.NotNil(h, i)
		assert.Equal(0, h.Len(), i)
		assert.Nil(h.Peek())
		assert.Nil(h.PopHeap())

		minVal := math.MaxFloat64
		for j, val := range tc.entries {
			if val.val < minVal {
				minVal = val.val
			}
			h.PushHeap(val)
			assert.Equal(j+1, h.Len(), i)
			assert.Equal(&testEntry{minVal}, h.Peek(), i)
			assert.Equal(j+1, h.Len(), i)
		}

		assert.Equal(5, h.Len())

		for i := 1; i <= 5; i++ {
			assert.Equal(&testEntry{float64(i)}, h.Peek())
			assert.Equal(&testEntry{float64(i)}, h.PopHeap())
			assert.Equal(5-i, h.Len())
		}

		assert.Equal(0, h.Len(), i)
		assert.Nil(h.Peek())
		assert.Nil(h.PopHeap())
	}
}

func TestAnyMatch(t *testing.T) {
	assert := assert.New(t)

	h1 := model.NewHeap(getVal)
	h1.PushHeap(&testEntry{2.34})
	h1.PushHeap(&testEntry{3.45})
	h1.PushHeap(&testEntry{1.23})
	h1.PushHeap(&testEntry{4})
	h1.PushHeap(&testEntry{5})
	h1.PushHeap(&testEntry{6})
	h1.PushHeap(&testEntry{7})

	assert.True(h1.AnyMatch(func(x interface{}) bool {
		e, okay := x.(*testEntry)
		return okay && e.val > 6
	}))

	assert.False(h1.AnyMatch(func(x interface{}) bool {
		_, okay := x.(*model2d.Circuit2D)
		return okay
	}))

	assert.False(h1.AnyMatch(func(x interface{}) bool {
		e, okay := x.(*testEntry)
		return okay && e.val < 1
	}))

	h2 := model.NewHeap(getVal)
	assert.False(h2.AnyMatch(func(x interface{}) bool {
		return true
	}))
}

func TestClone(t *testing.T) {
	assert := assert.New(t)

	h1 := model.NewHeap(getVal)
	h1.PushHeap(&testEntry{2.34})
	h1.PushHeap(&testEntry{3.45})
	h1.PushHeap(&testEntry{1.23})

	h2 := h1.Clone()
	assert.Equal(3, h2.Len())

	h1.PushHeap(&testEntry{1.5})
	assert.Equal(4, h1.Len())
	assert.Equal(3, h2.Len())

	assert.Equal(&testEntry{1.23}, h2.PopHeap())
	assert.Equal(4, h1.Len())
	assert.Equal(2, h2.Len())
	assert.Equal(&testEntry{2.34}, h2.PopHeap())
	assert.Equal(4, h1.Len())
	assert.Equal(1, h2.Len())

	h2.PushHeap(&testEntry{4.56})
	assert.Equal(4, h1.Len())
	assert.Equal(2, h2.Len())

	assert.Equal(&testEntry{1.23}, h1.PopHeap())
	assert.Equal(3, h1.Len())
	assert.Equal(2, h2.Len())
}

func TestDeleteAll(t *testing.T) {
	assert := assert.New(t)

	h1 := model.NewHeap(getVal)
	h1.PushHeap(&testEntry{2.34})
	h1.PushHeap(&testEntry{3.45})
	h1.PushHeap(&testEntry{1.23})
	h1.PushHeap(&testEntry{4})
	h1.PushHeap(&testEntry{5})
	h1.PushHeap(&testEntry{6})
	h1.PushHeap(&testEntry{7})

	h2 := h1.Clone()
	h2.DeleteAll(func(x interface{}) bool {
		return x.(*testEntry).val > 5.0
	})
	assert.Equal(7, h1.Len())
	assert.Equal(5, h2.Len())

	h3 := h1.Clone()
	h3.DeleteAll(func(x interface{}) bool {
		return x.(*testEntry).val >= 4.0
	})
	assert.Equal(7, h1.Len())
	assert.Equal(3, h3.Len())

	h2.DeleteAll(func(x interface{}) bool {
		return x.(*testEntry).val <= 2.5
	})
	assert.Equal(3, h2.Len())

	assert.Equal(&testEntry{3.45}, h2.PopHeap())
	assert.Equal(&testEntry{4}, h2.PopHeap())
	assert.Equal(&testEntry{5}, h2.PopHeap())
	assert.Equal(7, h1.Len())
	assert.Equal(0, h2.Len())
	assert.Equal(3, h3.Len())
}

func TestReplaceAll(t *testing.T) {
	assert := assert.New(t)

	h1 := model.NewHeap(getVal)
	h1.PushHeap(&testEntry{2.34})
	h1.PushHeap(&testEntry{3.45})
	h1.PushHeap(&testEntry{1.23})
	h1.PushHeap(&testEntry{4})
	h1.PushHeap(&testEntry{5})
	h1.PushHeap(&testEntry{6})
	h1.PushHeap(&testEntry{7})

	h2 := h1.Clone()
	h2.ReplaceAll(func(x interface{}) []interface{} {
		i := int(x.(*testEntry).val)
		if i%3 == 0 {
			return []interface{}{}
		} else if i%2 == 0 {
			return []interface{}{x, &testEntry{val: float64(i) * 2.0}, &testEntry{val: float64(i) * 0.5}}
		}
		return []interface{}{x}
	})

	assert.Equal(7, h1.Len())
	assert.Equal(9, h2.Len())

	assert.Equal(&testEntry{1.23}, h1.Peek())
	assert.Equal(&testEntry{1.23}, h1.PopHeap())
	assert.Equal(6, h1.Len())
	assert.Equal(9, h2.Len())

	assert.Equal(&testEntry{1.0}, h2.Peek())
	assert.Equal(&testEntry{1.0}, h2.PopHeap())
	assert.Equal(6, h1.Len())
	assert.Equal(8, h2.Len())

	assert.Equal(&testEntry{1.23}, h2.Peek())
	assert.Equal(&testEntry{1.23}, h2.PopHeap())
	assert.Equal(6, h1.Len())
	assert.Equal(7, h2.Len())

	assert.Equal(&testEntry{2.34}, h1.Peek())
	assert.Equal(&testEntry{2.34}, h1.PopHeap())
	assert.Equal(5, h1.Len())
	assert.Equal(7, h2.Len())

	assert.Equal(&testEntry{3.45}, h1.Peek())
	assert.Equal(&testEntry{3.45}, h1.PopHeap())
	assert.Equal(4, h1.Len())
	assert.Equal(7, h2.Len())

	assert.Equal(&testEntry{2.0}, h2.Peek())
	assert.Equal(&testEntry{2.0}, h2.PopHeap())
	assert.Equal(4, h1.Len())
	assert.Equal(6, h2.Len())

	assert.Equal(&testEntry{2.34}, h2.Peek())
	assert.Equal(&testEntry{2.34}, h2.PopHeap())
	assert.Equal(4, h1.Len())
	assert.Equal(5, h2.Len())

	assert.Equal(&testEntry{4.0}, h2.Peek())
	assert.Equal(&testEntry{4.0}, h2.PopHeap())
	assert.Equal(4, h1.Len())
	assert.Equal(4, h2.Len())

	assert.Equal(&testEntry{4.0}, h2.Peek())
	assert.Equal(&testEntry{4.0}, h2.PopHeap())
	assert.Equal(4, h1.Len())
	assert.Equal(3, h2.Len())
}

func TestToString_Heap(t *testing.T) {
	assert := assert.New(t)

	h1 := model.NewHeap(getVal)
	h1.PushHeap(&testEntry{2.34})
	h1.PushHeap(&testEntry{3.45})
	h1.PushHeap(&testEntry{1.23})
	h1.PushHeap(&testEntry{4})
	h1.PushHeap(&testEntry{5})
	h1.PushHeap(&testEntry{6})
	h1.PushHeap(&testEntry{7})

	assert.Equal(`{1.23,3.45,2.34,4,5,6,7}`, h1.ToString())
}
