package cache

import (
	"reflect"
	"sync"
	"time"
)

// Memoize returns a new pull through cache.
func Memoize(action func(interface{}) (interface{}, error), options ...MemoizeOption) *Memoized {
	m := Memoized{
		Action: action,
	}
	for _, opt := range options {
		opt(&m)
	}
	return &m
}

// MemoizeOption is an option for memoized functions.
type MemoizeOption func(*Memoized)

// OptMemoizeTTL sets the memoization ttl.
func OptMemoizeTTL(d time.Duration) MemoizeOption {
	return func(m *Memoized) {
		m.TTL = d
	}
}

// Memoized is a cache that memoizes the result of a function.
type Memoized struct {
	sync.Mutex
	TTL    time.Duration
	Action func(interface{}) (interface{}, error)
	Values []*MemoizedValue
}

// Call returns a cached response or calls the action with a given argument if it's expired.
// It uses a list and iterates over the values in the list to find the value.
// As a result, it sucks for high cardinality
func (m *Memoized) Call(args interface{}) (interface{}, error) {
	if args == nil {
		panic("nil args")
	}

	if !reflect.TypeOf(args).Comparable() {
		panic("args is not comparable")
	}

	m.Lock()
	defer m.Unlock()

	now := time.Now().UTC()
	for _, value := range m.Values {
		// if we have a matching args in the set
		if value.Args == args {
			// if there isn't a ttl, or we we have a ttl and the value isn't expired ...
			delta := now.Sub(value.Timestamp)
			if m.TTL == 0 || (m.TTL != 0 && delta < m.TTL) {
				return value.Response, value.Err
			}

			// refresh the value
			value.Response, value.Err = m.Action(args)
			value.Timestamp = time.Now().UTC()
			return value.Response, value.Err
		}
	}

	// add value
	val := MemoizedValue{
		Timestamp: time.Now().UTC(),
		Args:      args,
	}

	val.Response, val.Err = m.Action(args)
	m.Values = append(m.Values, &val)
	return val.Response, val.Err
}

// MemoizedValue is a specific invocation with an argument.
type MemoizedValue struct {
	Timestamp time.Time
	Args      interface{}
	Response  interface{}
	Err       error
}
