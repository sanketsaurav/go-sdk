package cache

import "container/heap"

var (
	emptyArray = make([]*Value, 0)
)

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

// Enqueue adds an element to the "back" of the queue.
func (lrh *LRUHeap) Enqueue(object *Value) {
	heap.Push(&lrh.Values, object)
}

// Dequeue removes the first (oldest) element from the queue.
func (lrh *LRUHeap) Dequeue() *Value {
	return heap.Pop(&lrh.Values).(*Value)
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
