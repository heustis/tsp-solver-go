package tspmodel_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/stretchr/testify/assert"
)

type getValObject interface {
	getVal() float64
}

type testEntry struct {
	val float64
}

func (e *testEntry) getVal() float64 {
	return e.val
}

func (e *testEntry) String() string {
	return fmt.Sprintf("%v", e.val)
}

type deletabletestEntry struct {
	testEntry
	isDeleted bool
}

func (e *deletabletestEntry) Delete() {
	e.isDeleted = true
}

func getVal(e interface{}) float64 {
	if v, okay := e.(getValObject); okay {
		return v.getVal()
	} else {
		return 0.0
	}
}

func BenchmarkReplaceAll(b *testing.B) {
	maxEntries := 10000000
	entries := make([]interface{}, maxEntries)
	for i := 0; i < maxEntries; i++ {
		entries[i] = &testEntry{val: float64(i)}
	}

	heap2 := tspmodel.NewHeap(getVal)

	replacementFunction := func(x interface{}) interface{} {
		i := int(x.(*testEntry).val)
		if i%3 == 0 {
			return nil
		} else if i%4 == 0 {
			return []interface{}{x, &testEntry{val: float64(i) * 0.5}, &testEntry{val: float64(i) * 1.5}, &testEntry{val: float64(i) * 3.0}}
		}
		return x
	}

	// BenchmarkReplaceAll/ReplaceAll-100000-16			169843116	     6.707 ns/op			2 B/op	       0 allocs/op
	// BenchmarkReplaceAll/ReplaceAll-1000000-16			10172	    130100 ns/op		60851 B/op	     933 allocs/op
	// BenchmarkReplaceAll/ReplaceAll-10000000-16				1	1021660300 ns/op	706686592 B/op	 8333333 allocs/op

	b.Run("ReplaceAll-100000", func(b *testing.B) {
		heap2.PushAll(entries[0:100000]...)
		heap2.ReplaceAll(replacementFunction)
	})
	heap2.Delete()
	heap2 = tspmodel.NewHeap(getVal)

	heap2.PushAll(entries[0:1000000]...)
	b.Run("ReplaceAll-1000000", func(b *testing.B) {
		heap2.ReplaceAll(replacementFunction)
	})
	heap2.Delete()
	heap2 = tspmodel.NewHeap(getVal)

	b.Run("ReplaceAll-10000000", func(b *testing.B) {
		heap2.PushAll(entries[0:10000000]...)
		heap2.ReplaceAll(replacementFunction)
	})
	heap2.Delete()
	heap2 = tspmodel.NewHeap(getVal)
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
		h := tspmodel.NewHeap(getVal)
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

	h1 := tspmodel.NewHeap(getVal)
	h1.PushHeap(&testEntry{2.34})
	h1.PushHeap(&testEntry{3.45})
	h1.PushHeap(&testEntry{1.23})
	h1.PushHeap(&testEntry{4})
	h1.PushHeap(&testEntry{5})
	h1.PushHeap(&testEntry{6})
	h1.PushHeap(&testEntry{7})

	assert.Equal([]interface{}{
		&testEntry{1.23}, &testEntry{3.45}, &testEntry{2.34}, &testEntry{4}, &testEntry{5}, &testEntry{6}, &testEntry{7},
	}, h1.GetValues())

	assert.True(h1.AnyMatch(func(x interface{}) bool {
		e, okay := x.(*testEntry)
		return okay && e.val > 6
	}))

	assert.False(h1.AnyMatch(func(x interface{}) bool {
		_, okay := x.(*circuit.ConvexConcave)
		return okay
	}))

	assert.False(h1.AnyMatch(func(x interface{}) bool {
		e, okay := x.(*testEntry)
		return okay && e.val < 1
	}))

	h2 := tspmodel.NewHeap(getVal)
	assert.False(h2.AnyMatch(func(x interface{}) bool {
		return true
	}))
}

func TestClone(t *testing.T) {
	assert := assert.New(t)

	h1 := tspmodel.NewHeap(getVal)
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

	h1 := tspmodel.NewHeap(getVal)
	h1.PushHeap(&testEntry{2.34})
	h1.PushHeap(&testEntry{3.45})
	h1.PushHeap(&testEntry{1.23})
	h1.PushHeap(&testEntry{4})
	h1.PushHeap(&testEntry{5})
	h1.PushHeap(&testEntry{6})

	deletableEntry := &deletabletestEntry{
		testEntry: testEntry{val: 7},
		isDeleted: false,
	}
	h1.PushHeap(deletableEntry)

	h2 := h1.Clone()
	h2.DeleteAll(func(x interface{}) bool {
		return x.(getValObject).getVal() > 5.0
	})
	assert.Equal(7, h1.Len())
	assert.Equal(5, h2.Len())
	assert.False(deletableEntry.isDeleted)

	h3 := h1.Clone()
	h3.DeleteAll(func(x interface{}) bool {
		return x.(getValObject).getVal() >= 4.0
	})
	assert.Equal(7, h1.Len())
	assert.Equal(3, h3.Len())

	h2.DeleteAll(func(x interface{}) bool {
		return x.(getValObject).getVal() <= 2.5
	})
	assert.Equal(3, h2.Len())

	assert.Equal(&testEntry{3.45}, h2.PopHeap())
	assert.Equal(&testEntry{4}, h2.PopHeap())
	assert.Equal(&testEntry{5}, h2.PopHeap())
	assert.Equal(7, h1.Len())
	assert.Equal(0, h2.Len())
	assert.Equal(3, h3.Len())

	assert.False(deletableEntry.isDeleted)
	h1.Delete()
	assert.True(deletableEntry.isDeleted)
}

func TestHeapify(t *testing.T) {
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
		h := tspmodel.NewHeap(getVal)

		for j, val := range tc.entries {
			h.Push(val)
			assert.Equal(j+1, h.Len(), i)
			assert.Equal(tc.entries[0], h.Peek(), i)
		}

		h.Heapify()

		for i := 1; i <= 5; i++ {
			assert.Equal(&testEntry{float64(i)}, h.PopHeap())
			assert.Equal(5-i, h.Len())
		}
	}
}

func TestPushAll(t *testing.T) {
	assert := assert.New(t)

	entries := []interface{}{&testEntry{5}, &testEntry{4}, &testEntry{3}, &testEntry{2}, &testEntry{1}}

	h := tspmodel.NewHeap(getVal)
	h.PushAll(entries...)
	assert.NotEqual(entries, h.GetValues())
	assert.Equal([]interface{}{&testEntry{1}, &testEntry{2}, &testEntry{3}, &testEntry{5}, &testEntry{4}}, h.GetValues())

	h2 := tspmodel.NewHeap(getVal)
	h2.PushAllFrom(h)
	assert.Equal([]interface{}{&testEntry{1}, &testEntry{2}, &testEntry{3}, &testEntry{5}, &testEntry{4}}, h2.GetValues())

	h2.PopHeap()
	assert.Equal([]interface{}{&testEntry{2}, &testEntry{4}, &testEntry{3}, &testEntry{5}}, h2.GetValues())
	assert.Equal([]interface{}{&testEntry{1}, &testEntry{2}, &testEntry{3}, &testEntry{5}, &testEntry{4}}, h.GetValues())

	h.PopHeap()
	h.PopHeap()
	assert.Equal([]interface{}{&testEntry{2}, &testEntry{4}, &testEntry{3}, &testEntry{5}}, h2.GetValues())
	assert.Equal([]interface{}{&testEntry{3}, &testEntry{4}, &testEntry{5}}, h.GetValues())

	h.Delete()
	assert.Equal(0, h.Len())
	h2.Delete()
	assert.Equal(0, h2.Len())
}

func TestReplaceAll(t *testing.T) {
	assert := assert.New(t)

	h1 := tspmodel.NewHeap(getVal)
	h1.PushHeap(&testEntry{2.34})
	h1.PushHeap(&testEntry{3.45})
	h1.PushHeap(&testEntry{1.23})
	h1.PushHeap(&testEntry{4})
	h1.PushHeap(&testEntry{5})
	h1.PushHeap(&testEntry{6})
	h1.PushHeap(&testEntry{7})

	h2 := h1.Clone()
	h2.ReplaceAll(func(x interface{}) interface{} {
		i := int(x.(*testEntry).val)
		if i%3 == 0 {
			return nil
		} else if i%2 == 0 {
			return []interface{}{x, &testEntry{val: float64(i) * 2.0}, &testEntry{val: float64(i) * 0.5}}
		}
		return x
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

func TestTrimN(t *testing.T) {
	assert := assert.New(t)

	entries := []interface{}{&testEntry{5}, &testEntry{4}, &testEntry{3}, &testEntry{2}, &testEntry{1}}

	h := tspmodel.NewHeap(getVal)
	h.PushAll(entries...)
	assert.Equal(5, h.Len())

	h.TrimN(4)
	assert.Equal(4, h.Len())
	assert.Equal([]interface{}{&testEntry{1}, &testEntry{2}, &testEntry{3}, &testEntry{4}}, h.GetValues())

	h.TrimN(5)
	assert.Equal(4, h.Len())
	assert.Equal([]interface{}{&testEntry{1}, &testEntry{2}, &testEntry{3}, &testEntry{4}}, h.GetValues())

	h.TrimN(2)
	assert.Equal(2, h.Len())
	assert.Equal([]interface{}{&testEntry{1}, &testEntry{2}}, h.GetValues())

	h.TrimN(-1)
	assert.Equal(2, h.Len())
	assert.Equal([]interface{}{&testEntry{1}, &testEntry{2}}, h.GetValues())

	h.TrimN(0)
	assert.Equal(0, h.Len())
}

func TestString_Heap(t *testing.T) {
	assert := assert.New(t)

	h1 := tspmodel.NewHeap(getVal)
	h1.PushHeap(&testEntry{2.34})
	h1.PushHeap(&testEntry{3.45})
	h1.PushHeap(&testEntry{1.23})
	h1.PushHeap(&testEntry{4})
	h1.PushHeap(&testEntry{5})
	h1.PushHeap(&testEntry{6})
	h1.PushHeap(&testEntry{7})

	assert.Equal(`{1.23,3.45,2.34,4,5,6,7}`, h1.String())
}

func TestString_Heap2(t *testing.T) {
	assert := assert.New(t)

	h1 := tspmodel.NewHeap(getVal)
	h1.PushHeap(1.234)
	h1.PushHeap("test")

	d := &tspmodel.DistanceToEdge{
		Vertex:   tspmodel2d.NewVertex2D(123.45, 678.9),
		Edge:     tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(5.15, 0.13), tspmodel2d.NewVertex2D(1000.3, 1100.25)),
		Distance: 567.89,
	}
	h1.PushHeap(d)

	type testStruct struct {
		Foo   string `json:"bar"`
		Other int    `json:"other"`
	}

	h1.PushHeap(&testStruct{Foo: "test data", Other: 567})

	assert.Equal(`{1.234,"test",{"vertex":{"x":123.45,"y":678.9},"edge":{"start":{"x":5.15,"y":0.13},"end":{"x":1000.3,"y":1100.25}},"distance":567.89},{"bar":"test data","other":567}}`, h1.String())
}
