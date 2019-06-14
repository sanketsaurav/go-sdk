package cache

import "container/heap"

// NewLRUHeap creates a new, empty, LRU Heap.
func NewLRUHeap() *LRUHeap {
	return &LRUHeap{}
}

// LRUHeap is a fifo buffer that is backed by a pre-allocated array, instead of allocating
// a whole new node object for each element (which saves GC churn).
// Enqueue can be O(n), Dequeue can be O(1).
type LRUHeap struct {
	Values lruHeapValues
}

// Len returns the length of the queue (as it is currently populated).
// Actual memory footprint may be different.
func (lrh *LRUHeap) Len() int {
	return len(lrh.Values)
}

// Clear removes all objects from the queue.
func (lrh *LRUHeap) Clear() {
	lrh.Values = nil
}

// Push adds an element to the heap.
func (lrh *LRUHeap) Push(object *Value) {
	heap.Push(&lrh.Values, object)
}

// Pop removes the first (oldest) element from the heap.
func (lrh *LRUHeap) Pop() *Value {
	if len(lrh.Values) == 0 {
		return nil
	}
	return heap.Pop(&lrh.Values).(*Value)
}

// Fix updates a value by key.
func (lrh *LRUHeap) Fix(key interface{}, newValue *Value) {
	var i int
	for index, value := range lrh.Values {
		if value.Key == key {
			i = index
			break
		}
	}
	lrh.Values[i] = newValue
	heap.Fix(&lrh.Values, i)
}

// RemoveByKey removes a value by key.
func (lrh *LRUHeap) RemoveByKey(key interface{}) {
	var i int
	for index, value := range lrh.Values {
		if value.Key == key {
			i = index
			break
		}
	}
	heap.Remove(&lrh.Values, i)
}

// Peek returns the oldest value but does not dequeue it.
func (lrh *LRUHeap) Peek() *Value {
	if len(lrh.Values) == 0 {
		return nil
	}
	return lrh.Values[0]
}

// ConsumeUntil calls the consumer for each element in the buffer, while also dequeueing that entry.
// The consumer should return `true` if it should continue processing.
func (lrh *LRUHeap) ConsumeUntil(consumer func(value *Value) bool) {
	if len(lrh.Values) == 0 {
		return
	}

	len := len(lrh.Values)
	for i := 0; i < len; i++ {
		if !consumer(lrh.Peek()) {
			return
		}
		lrh.Pop()
	}
}

var (
	_ heap.Interface = (*lruHeapValues)(nil)
)

type lruHeapValues []*Value

func (lruv lruHeapValues) Len() int           { return len(lruv) }
func (lruv lruHeapValues) Less(i, j int) bool { return lruv[i].Timestamp.Before(lruv[j].Timestamp) }
func (lruv lruHeapValues) Swap(i, j int)      { lruv[i], lruv[j] = lruv[j], lruv[i] }

func (lruv *lruHeapValues) Push(x interface{}) {
	*lruv = append(*lruv, x.(*Value))
}

func (lruv *lruHeapValues) Pop() interface{} {
	old := *lruv
	n := len(old)
	x := old[n-1]
	*lruv = old[0 : n-1]
	return x
}
