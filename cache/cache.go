package cache

import "sync"

// Cache is a type that implements the cache interface.
type Cache interface {
	sync.Locker
	RLock()
	RUnlock()

	Has(key interface{}) bool
	Set(key, value interface{}, options ...ValueOption)
	Get(key interface{}) (interface{}, bool)
	Remove(key interface{}) (interface{}, bool)
}
