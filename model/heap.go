package model

import (
	"container/heap"
	"math"
)

// Heap is a min heap that implements heap.Interface.
// It provides the methods PushHeap and PopHeap, so that clients don't need to use the heap package.
// It also provides methods for cloning the heap, deleting items from the heap based on a condition function, and peeking at the head of the heap.
type Heap struct {
	arr   []interface{}
	value func(a interface{}) float64
}

// NewHeap creates a new, empty, min Heap with the supplied function for computing the value used to order the heap.
func NewHeap(valueFunc func(a interface{}) float64) *Heap {
	return &Heap{
		arr:   []interface{}{},
		value: valueFunc,
	}
}

// AnyMatch checks if any items in the heap match the supplied predicate, and returns true if any item does.
// This method will return true as soon as any matching item is found.
// If the heap is empty, this method will return false.
// The complexity is O(n) due to needing to potentially check all the items in the heap.
func (h *Heap) AnyMatch(predicate func(x interface{}) bool) bool {
	for _, item := range h.arr {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Clone creates a copy of the Heap, so that items can be added or removed from one heap without affecting the other heap.
// The complexity is O(n) due to needing to copy all the items in the heap.
func (h *Heap) Clone() *Heap {
	arrCopy := make([]interface{}, len(h.arr))
	copy(arrCopy, h.arr)
	return &Heap{
		arr:   arrCopy,
		value: h.value,
	}
}

// Delete cleans up this heap by: unsetting the value function (since it may reference the object containing the heap),
// deleting all entries in the heap that implement Deletable, and clearing the backing array.
func (h *Heap) Delete() {
	h.value = nil
	for i, entry := range h.arr {
		if del, okay := entry.(Deletable); okay {
			del.Delete()
		}
		h.arr[i] = nil
	}
	h.arr = []interface{}{}
}

// DeleteAll creates a copy of the heap, only containing items that match the supplied condition; this modifies the original heap.
// This returns all items that were deleted in an array.
// "shouldDelete" needs to return true if the item should be deleted, false if it should be included in the copy.
// The complexity is O(n) due to needing to check each entry for inclusion, then re-heapifying the remaining data.
func (h *Heap) DeleteAll(shouldDelete func(x interface{}) bool) []interface{} {
	updated := make([]interface{}, 0, h.Len())
	deleted := []interface{}{}

	for _, v := range h.arr {
		if shouldDelete(v) {
			deleted = append(deleted, v)
		} else {
			updated = append(updated, v)
		}
	}

	h.arr = updated
	heap.Init(h)

	return deleted
}

// GetValues returns all items in the heap. If any of the items in the heap (or the array) are modified, Heapify() should be called on this heap.
func (h *Heap) GetValues() []interface{} {
	return h.arr
}

// Heapify ensures that the heap is ordered correctly.
// The complexity is O(n) per heap.Init().
func (h *Heap) Heapify() {
	heap.Init(h)
}

// Push is required by heap.Interface, and is used by the heap package.
// PushHeap or PushAll should be used instead of Push to add a new item to the heap, since this will not update the location of the new item to the correct position in the heap.
func (h *Heap) Push(x interface{}) {
	h.arr = append(h.arr, x)
}

// PushAll adds all the supplied elements to the heap.
// The complexity is O(n), where n is the number elements already in the heap + the number of elements in the supplied array.
func (h *Heap) PushAll(x ...interface{}) {
	h.arr = append(h.arr, x...)
	heap.Init(h)
}

// PushAllFrom adds all the entries in the supplied heap to this heap.
// The complexity is O(n), where n is the number elements already in the heap + the number of elements in the supplied array.
func (h *Heap) PushAllFrom(other *Heap) {
	h.PushAll(other.arr...)
}

// PushHeap adds the supplied item to the heap at the correct location.
// This wraps heap.Push so that clients don't need to be familiar with the heap package, and only need to use this class.
// The complexity is O(log n).
func (h *Heap) PushHeap(x interface{}) {
	heap.Push(h, x)
}

// Pop is required by heap.Interface, and is used by the heap package.
// PopHeap should be used instead of Pop to remove an item from the heap, since this will not re-order the remaining heap.
func (h *Heap) Pop() interface{} {
	lastIndex := h.Len() - 1
	if lastIndex < 0 {
		return nil
	}
	toReturn := h.arr[lastIndex]
	h.arr = h.arr[:lastIndex]
	return toReturn
}

// PopHeap removes the next item from the heap (based the minimum value of any item when the configured value function is applied).
// This wraps heap.Pop so that clients don't need to be familiar with the heap package, and only need to use this class.
// The complexity is O(log n).
func (h *Heap) PopHeap() interface{} {
	return heap.Pop(h)
}

// Peek retrieves the minimum value from the heap, without modifying the heap. The complexity is O(1).
func (h *Heap) Peek() interface{} {
	if h.Len() <= 0 {
		return nil
	}
	return h.arr[0]
}

// Len retrieves the number of items, without modifying the heap. The complexity is O(1).
func (h *Heap) Len() int {
	return len(h.arr)
}

// Less is required by heap.Interface and is used to compare items in the heap when sorting the heap. The complexity is O(1).
func (h *Heap) Less(i int, j int) bool {
	return h.value(h.arr[i]) < h.value(h.arr[j])
}

// ReplaceAll creates a copy of the heap, using the supplied function to replace items in the heap; this modifies the original heap.
// "replaceFunction" should return one of the following:
//   - an empty array if the item should be excluded,
//   - an array containing the original item if it should be retained unchanged, or
//   - an array with one or more items if the original item should be replaced.
// The complexity is O(n) due to needing to check each entry for replacement, then re-heapifying the updated data.
func (h *Heap) ReplaceAll(replaceFunction func(x interface{}) []interface{}) {
	updated := make([]interface{}, 0, h.Len())

	for _, v := range h.arr {
		updated = append(updated, replaceFunction(v)...)
	}

	h.arr = updated
	heap.Init(h)
}

// ReplaceAll2 creates a copy of the heap, using the supplied function to replace items in the heap; this modifies the original heap.
// "replaceFunction" should return one of the following:
//   - nil or an empty array if the item should be excluded,
//   - a single item that should replace the original item,
//   - the original item if it should be retained unchanged, or
//   - an array with one or more items if the original item should be replaced.
// The complexity is O(n) due to needing to check each entry for replacement, then re-heapifying the updated data.
func (h *Heap) ReplaceAll2(replaceFunction func(x interface{}) interface{}) {
	updated := make([]interface{}, 0, h.Len())

	for _, x := range h.arr {
		replacement := replaceFunction(x)
		if replacementArray, okay := replacement.([]interface{}); okay {
			updated = append(updated, replacementArray...)
		} else if replacement != nil {
			updated = append(updated, replacement)
		}
	}

	h.arr = updated
	heap.Init(h)
}

// Swap is required by heap.Interface and is used swap items in the heap when sorting the heap. The complexity is O(1).
func (h *Heap) Swap(i int, j int) {
	if i >= 0 && j >= 0 {
		h.arr[i], h.arr[j] = h.arr[j], h.arr[i]
	}
}

// TrimN keeps the N minmimum entries in this heap and discards the rest. The complexity is O(numberToRetain * log(n)).
func (h *Heap) TrimN(numberToRetain int) {
	numberToRetain = int(math.Min(float64(numberToRetain), float64(h.Len())))
	updated := make([]interface{}, numberToRetain)

	for i := 0; i < numberToRetain; i++ {
		updated[i] = h.PopHeap()
	}
	h.arr = updated
	heap.Init(h)
}

func (h *Heap) ToString() string {
	str := "{"

	for i, entry := range h.arr {
		if i == 0 {
			str += ToString(entry)
		} else {
			str += "," + ToString(entry)
		}
	}
	return str + "}"
}

var _ heap.Interface = (*Heap)(nil)
var _ Deletable = (*Heap)(nil)
