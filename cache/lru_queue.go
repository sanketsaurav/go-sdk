package cache

const (
	ringBufferMinimumGrow     = 4
	ringBufferShrinkThreshold = 32
	ringBufferGrowFactor      = 200
	ringBufferDefaultCapacity = 4
)

var (
	emptyArray = make([]interface{}, 0)
)

// NewLRUQueue creates a new, empty, RingBuffer.
func NewLRUQueue() *LRUQueue {
	return &LRUQueue{
		array: make([]*Value, ringBufferDefaultCapacity),
		head:  0,
		tail:  0,
		size:  0,
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

// Len returns the length of the queue (as it is currently populated).
// Actual memory footprint may be different.
func (lrq *LRUQueue) Len() (len int) {
	return lrq.size
}

// Capacity returns the total size of the queue, including empty elements.
func (lrq *LRUQueue) Capacity() int {
	return len(lrq.array)
}

// Clear removes all objects from the queue.
func (lrq *LRUQueue) Clear() {
	if lrq.head < lrq.tail {
		arrayClear(lrq.array, lrq.head, lrq.size)
	} else {
		arrayClear(lrq.array, lrq.head, len(lrq.array)-lrq.head)
		arrayClear(lrq.array, 0, lrq.tail)
	}

	lrq.head = 0
	lrq.tail = 0
	lrq.size = 0
}

// Enqueue adds an element to the "back" of the queue.
func (lrq *LRUQueue) Enqueue(object *Value) {
	if lrq.size == len(lrq.array) {
		newCapacity := int(len(lrq.array) * int(ringBufferGrowFactor/100))
		if newCapacity < (len(lrq.array) + ringBufferMinimumGrow) {
			newCapacity = len(lrq.array) + ringBufferMinimumGrow
		}
		lrq.setCapacity(newCapacity)
	}

	lrq.array[lrq.tail] = object
	lrq.tail = (lrq.tail + 1) % len(lrq.array)
	lrq.size++
}

// Dequeue removes the first (oldest) element from the queue.
func (lrq *LRUQueue) Dequeue() interface{} {
	if lrq.size == 0 {
		return nil
	}

	removed := lrq.array[lrq.head]
	lrq.head = (lrq.head + 1) % len(lrq.array)
	lrq.size--
	return removed
}

// Peek returns but does not remove the first element.
func (lrq *LRUQueue) Peek() interface{} {
	if lrq.size == 0 {
		return nil
	}
	return lrq.array[lrq.head]
}

// PeekBack returns but does not remove the last element.
func (lrq *LRUQueue) PeekBack() interface{} {
	if lrq.size == 0 {
		return nil
	}
	if lrq.tail == 0 {
		return lrq.array[len(lrq.array)-1]
	}
	return lrq.array[lrq.tail-1]
}

func (lrq *LRUQueue) setCapacity(capacity int) {
	newArray := make([]*Value, capacity)
	if lrq.size > 0 {
		if lrq.head < lrq.tail {
			arrayCopy(lrq.array, lrq.head, newArray, 0, lrq.size)
		} else {
			arrayCopy(lrq.array, lrq.head, newArray, 0, len(lrq.array)-lrq.head)
			arrayCopy(lrq.array, 0, newArray, len(lrq.array)-lrq.head, lrq.tail)
		}
	}
	lrq.array = newArray
	lrq.head = 0
	if lrq.size == capacity {
		lrq.tail = 0
	} else {
		lrq.tail = lrq.size
	}
}

// trimExcess resizes the buffer to better fit the contents.
func (lrq *LRUQueue) trimExcess() {
	threshold := float64(len(lrq.array)) * 0.9
	if lrq.size < int(threshold) {
		lrq.setCapacity(lrq.size)
	}
}

// Each calls the consumer for each element in the buffer.
func (lrq *LRUQueue) Each(consumer func(value interface{})) {
	if lrq.size == 0 {
		return
	}

	if lrq.head < lrq.tail {
		for cursor := lrq.head; cursor < lrq.tail; cursor++ {
			consumer(lrq.array[cursor])
		}
	} else {
		for cursor := lrq.head; cursor < len(lrq.array); cursor++ {
			consumer(lrq.array[cursor])
		}
		for cursor := 0; cursor < lrq.tail; cursor++ {
			consumer(lrq.array[cursor])
		}
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
