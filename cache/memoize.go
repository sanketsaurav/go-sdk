package cache

import (
	"reflect"
	"sync"
	"time"
)

// Memoize returns a new pull through cache.
func Memoize(action func(interface{}) (interface{}, error), options ...MemoizeOption) *Memoized {
	return &Memoized{
		Action: action,
		TTL:    5 * time.Second,
	}
}

// MemoizeOption is an option for memoized functions.
type MemoizeOption func(*Memoized)

// Memoized is a cache that memoizes the result of a function.
type Memoized struct {
	sync.RWMutex
	TTL    time.Duration
	Action func(interface{}) (interface{}, error)
	Values []*MemoizedValue
}

// Call returns a cached response or calls the action with a given argument if it's expired.
func (m *Memoized) Call(args interface{}) (interface{}, error) {
	if args == nil {
		panic("nil args")
	}

	if !reflect.TypeOf(args).Comparable() {
		panic("args is not comparable")
	}

	m.RLock()
	for _, value := range m.Values {
		if value.Args == args {
			m.RUnlock()
			return value.Call(args, m.Action)
		}
	}
	return nil, nil
}

// MemoizedValue is a specific invocation with an argument.
type MemoizedValue struct {
	sync.Mutex

	Args      interface{}
	Response  interface{}
	Err       error
	Timestamp time.Duration
	TTL       time.Duration
}

// Call returns either a cached value or a new call.
func (m *MemoizedValue) Call(args interface{}, action func(interface{}) (interface{}, error)) (interface{}, error) {
	return nil, nil
}
