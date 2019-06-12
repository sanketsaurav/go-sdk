package cache

import (
	"time"
)

// ValueOption is an option for a cache value.
type ValueOption func(*Value)

// OptValueTTL sets the ttl for the value.
func OptValueTTL(d time.Duration) ValueOption {
	return func(v *Value) {
		v.TTL = d
	}
}

// OptValueTimestamp sets the timestamp for the value.
func OptValueTimestamp(t time.Time) ValueOption {
	return func(v *Value) {
		v.Timestamp = t
	}
}

// OptValueOnSweep sets the on sweep handler.
func OptValueOnSweep(handler func()) ValueOption {
	return func(v *Value) {
		v.OnSweep = handler
	}
}

// OptValueOnRemove sets the on remove handler.
func OptValueOnRemove(handler func()) ValueOption {
	return func(v *Value) {
		v.OnRemove = handler
	}
}

// Value is a cached item.
type Value struct {
	Timestamp time.Time
	Key       interface{}
	Value     interface{}
	TTL       time.Duration
	OnRemove  func()
	OnSweep   func()
}
