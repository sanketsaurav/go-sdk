package cache

// LRU is a type that implements the LRU methods.
type LRU interface {
	// Len returns the number of items in the queue.
	Len() int
	// Push should add a new value. The new minimum value should be returned by `Peek()`.
	Push(*Value)
	// Pop should remove and return the minimum value.
	Pop() *Value
	// Peek should return the minimum value.
	Peek() *Value
	// Fix should update the LRU, replacing any existing values.
	Fix(*Value)
	// Remove should remove a value with a given key.
	Remove(interface{})
	// Consume should iterate through the values. If `true` is removed by the handler,
	// the current value will be removed and the handler will be called on the next value.
	Consume(func(*Value) bool)
}
