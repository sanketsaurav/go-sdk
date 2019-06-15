package cache

import "sort"

var (
	_ LRU = (*LRUQueue)(nil)
)

// NewLRUQueue creates a new, empty, LRUQueue.
func NewLRUQueue() *LRUQueue {
	return &LRUQueue{
		array: make([]*Value, ringBufferDefaultCapacity),
	}
}

// LRUQueue is a fifo buffer that is backed by a pre-allocated array, instead of allocating
// a whole new node object for each element (which saves GC churn).
// Enqueue can be O(n), Dequeue can be O(1).
type LRUQueue struct {
	array []*Value
	head  int
	tail  int
	size  int
}

// Len returns the length of the ring buffer (as it is currently populated).
// Actual memory footprint may be different.
func (lru *LRUQueue) Len() (len int) {
	return lru.size
}

// Capacity returns the total size of the ring bufffer, including empty elements.
func (lru *LRUQueue) Capacity() int {
	return len(lru.array)
}

// Clear removes all objects from the LRUQueue.
func (lru *LRUQueue) Clear() {
	if lru.head < lru.tail {
		arrayClear(lru.array, lru.head, lru.size)
	} else {
		arrayClear(lru.array, lru.head, len(lru.array)-lru.head)
		arrayClear(lru.array, 0, lru.tail)
	}
	lru.head = 0
	lru.tail = 0
	lru.size = 0
}

// Push adds an element to the "back" of the LRUQueue.
func (lru *LRUQueue) Push(object *Value) {
	if lru.size == len(lru.array) {
		lru.setCapacity(lru.newCapacity())
	}
	lru.array[lru.tail] = object
	lru.tail = (lru.tail + 1) % len(lru.array)
	lru.size++
}

// Pop removes the first (oldest) element from the LRUQueue.
func (lru *LRUQueue) Pop() *Value {
	if lru.size == 0 {
		return nil
	}

	removed := lru.array[lru.head]
	lru.head = (lru.head + 1) % len(lru.array)
	lru.size--
	return removed
}

// Peek returns but does not remove the first element.
func (lru *LRUQueue) Peek() *Value {
	if lru.size == 0 {
		return nil
	}
	return lru.array[lru.head]
}

// PeekBack returns but does not remove the last element.
func (lru *LRUQueue) PeekBack() *Value {
	if lru.size == 0 {
		return nil
	}
	if lru.tail == 0 {
		return lru.array[len(lru.array)-1]
	}
	return lru.array[lru.tail-1]
}

// Fix updates the queue given an update to a specific value.
func (lru *LRUQueue) Fix(value *Value) {
	if lru.size == 0 {
		return
	}
	if value == nil {
		panic("lru queue; value is nil")
	}

	values := make([]*Value, lru.size)

	var cursorValue *Value
	var index int
	if lru.head < lru.tail {
		for cursor := lru.head; cursor < lru.tail; cursor++ {
			cursorValue = lru.array[cursor]
			if cursorValue.Key == value.Key {
				values[index] = value
			} else if cursorValue != nil {
				values[index] = cursorValue
			}
			index++
		}
	} else {
		for cursor := lru.head; cursor < len(lru.array); cursor++ {
			cursorValue = lru.array[cursor]
			if cursorValue.Key == value.Key {
				values[index] = value
			} else if cursorValue != nil {
				values[index] = cursorValue
			}
			index++
		}
		for cursor := 0; cursor < lru.tail; cursor++ {
			cursorValue = lru.array[cursor]
			if cursorValue.Key == value.Key {
				values[index] = value
			} else if cursorValue != nil {
				values[index] = cursorValue
			}
			index++
		}
	}

	// sort and recreate
	sort.Sort(LRUHeapValues(values))

	lru.array = values
	lru.head = 0
	lru.tail = len(values)
	lru.size = len(values)
}

// Remove removes an item from the queue by its key.
func (lru *LRUQueue) Remove(key interface{}) {
	if lru.size == 0 {
		return
	}

	if lru.size == 1 {
		if lru.array[lru.head].Key == key {
			lru.Pop()
		}
		return
	}

	// set up a new values list
	values := make([]*Value, lru.size-1)

	var cursorValue *Value
	var index int
	if lru.head < lru.tail {
		for cursor := lru.head; cursor < lru.tail; cursor++ {
			cursorValue = lru.array[cursor]
			if cursorValue.Key == key {
				continue
			}
			values[index] = cursorValue
			index++
		}
	} else {
		for cursor := lru.head; cursor < len(lru.array); cursor++ {
			cursorValue = lru.array[cursor]
			if cursorValue.Key == key {
				continue
			}
			values[index] = cursorValue
			index++
		}
		for cursor := 0; cursor < lru.tail; cursor++ {
			cursorValue = lru.array[cursor]
			if cursorValue.Key == key {
				continue
			}
			values[index] = cursorValue
			index++
		}
	}

	// recreate
	lru.array = values
	lru.head = 0
	lru.tail = len(values) - 1
	lru.size = len(values)
}

// ConsumeUntil calls the consumer for each element in the buffer. If the handler returns true,
// the element is popped and the handler is called on the next value.
func (lru *LRUQueue) ConsumeUntil(consumer func(value *Value) bool) {
	if lru.size == 0 {
		return
	}

	for i := 0; i < lru.size; i++ {
		if !consumer(lru.Peek()) {
			return
		}
		lru.Pop()
	}
}

//
// util / helpers
//

func (lru *LRUQueue) newCapacity() int {
	newCapacity := int(len(lru.array) * int(ringBufferGrowFactor/100))
	if newCapacity < (len(lru.array) + ringBufferMinimumGrow) {
		newCapacity = len(lru.array) + ringBufferMinimumGrow
	}
	return newCapacity
}

func (lru *LRUQueue) setCapacity(capacity int) {
	newArray := make([]*Value, capacity)
	if lru.size > 0 {
		if lru.head < lru.tail {
			arrayCopy(lru.array, lru.head, newArray, 0, lru.size)
		} else {
			arrayCopy(lru.array, lru.head, newArray, 0, len(lru.array)-lru.head)
			arrayCopy(lru.array, 0, newArray, len(lru.array)-lru.head, lru.tail)
		}
	}
	lru.array = newArray
	lru.head = 0
	if lru.size == capacity {
		lru.tail = 0
	} else {
		lru.tail = lru.size
	}
}

func (lru *LRUQueue) trimExcess() {
	threshold := float64(len(lru.array)) * 0.9
	if lru.size < int(threshold) {
		lru.setCapacity(lru.size)
	}
}

//
// array helpers
//

func arrayClear(source []*Value, index, length int) {
	for x := 0; x < length; x++ {
		absoluteIndex := x + index
		source[absoluteIndex] = nil
	}
}

func arrayCopy(source []*Value, sourceIndex int, destination []*Value, destinationIndex, length int) {
	for x := 0; x < length; x++ {
		from := sourceIndex + x
		to := destinationIndex + x

		destination[to] = source[from]
	}
}

const (
	ringBufferMinimumGrow     = 4
	ringBufferShrinkThreshold = 32
	ringBufferGrowFactor      = 200
	ringBufferDefaultCapacity = 4
)
