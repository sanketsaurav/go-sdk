package cache

const (
	ringBufferMinimumGrow     = 4
	ringBufferShrinkThreshold = 32
	ringBufferGrowFactor      = 200
	ringBufferDefaultCapacity = 4
)

// NewLRUQueue creates a new, empty, LRUQueue.
func NewLRUQueue(values ...*Value) *LRUQueue {
	return &LRUQueue{
		array: values,
		tail:  len(values) - 1,
		size:  len(values),
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

// Enqueue adds an element to the "back" of the LRUQueue.
func (lru *LRUQueue) Enqueue(object *Value) {
	if lru.size == len(lru.array) {
		newCapacity := int(len(lru.array) * int(ringBufferGrowFactor/100))
		if newCapacity < (len(lru.array) + ringBufferMinimumGrow) {
			newCapacity = len(lru.array) + ringBufferMinimumGrow
		}
		lru.setCapacity(newCapacity)
	}

	lru.array[lru.tail] = object
	lru.tail = (lru.tail + 1) % len(lru.array)
	lru.size++
}

// Dequeue removes the first (oldest) element from the LRUQueue.
func (lru *LRUQueue) Dequeue() *Value {
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

// ConsumeUntil calls the consumer for each element in the buffer, while also dequeueing that entry.
func (lru *LRUQueue) ConsumeUntil(consumer func(value *Value) bool) {
	if lru.size == 0 {
		return
	}

	len := lru.Len()
	for i := 0; i < len; i++ {
		if consumer(lru.Peek()) {
			lru.Dequeue()
			return
		}
		return
	}
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

// trimExcess resizes the buffer to better fit the contents.
func (lru *LRUQueue) trimExcess() {
	threshold := float64(len(lru.array)) * 0.9
	if lru.size < int(threshold) {
		lru.setCapacity(lru.size)
	}
}

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
