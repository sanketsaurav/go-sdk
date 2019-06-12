package cache

import (
	"reflect"
	"sync"
	"time"
)

// New returns a new cache.
func New(options ...Option) *Cache {
	c := Cache{
		Data:         make(map[interface{}]*Value),
		SweepStarted: make(chan struct{}),
	}
	for _, opt := range options {
		opt(&c)
	}
	return &c
}

// Option mutates the cache.
type Option func(*Cache)

// OptSweepInterval sets the cache interval.
func OptSweepInterval(interval time.Duration) Option {
	return func(c *Cache) {
		c.SweepInterval = interval
	}
}

// Cache is a memory cache.
type Cache struct {
	sync.RWMutex
	Data map[interface{}]*Value

	SweepInterval time.Duration

	SweepStarted  chan struct{}
	SweepCancel   chan struct{}
	SweepCanceled chan struct{}
}

// SweepIntervalOrDefault returns the sweep interval or a default.
func (c *Cache) SweepIntervalOrDefault() time.Duration {
	if c.SweepInterval > 0 {
		return c.SweepInterval
	}
	return 500 * time.Millisecond
}

// StartSweeping does the interval sweep for expired cache entries.
// You should call this yourself after you've created the cache object
//    c := cache.New()
//    go c.StartSweeping()
//    <-c.SweepStarted
//    ... code continues ...
func (c *Cache) StartSweeping() {
	if c.SweepStarted != nil {
		close(c.SweepStarted)
	}

	ticker := time.Tick(c.SweepIntervalOrDefault())
	c.SweepCancel = make(chan struct{})
	for {
		select {
		case <-c.SweepCancel:
			close(c.SweepCanceled)
			return
		default:
		}

		select {
		case <-c.SweepCancel:
			close(c.SweepCanceled)
			return
		case <-ticker:
			c.Sweep()
		}
	}
}

// Sweep checks keys for expired ttls.
// If any values are configured with 'OnSweep' handlers, they will be called
// outside holding the critical section.
func (c *Cache) Sweep() {
	c.Lock()
	now := time.Now().UTC()
	var keysToRemove []interface{}
	var handlers []func()
	for key, value := range c.Data {
		if now.After(value.Timestamp.Add(value.TTL)) {
			keysToRemove = append(keysToRemove, key)
			if value.OnSweep != nil {
				handlers = append(handlers, value.OnSweep)
			}
		}
	}
	for _, key := range keysToRemove {
		delete(c.Data, key)
	}
	c.Unlock()

	// call the handlers outside the critical section.
	for _, handler := range handlers {
		handler()
	}
}

// Set adds a cache item.
func (c *Cache) Set(key, value interface{}, options ...ValueOption) {
	if key == nil {
		panic("nil key")
	}

	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	if reflect.TypeOf(key).Kind() != reflect.Struct {
		panic("key is not a struct; consider using a key type to leverage the compiler for key checking")
	}

	v := Value{
		Timestamp: time.Now().UTC(),
		Key:       key,
		Value:     value,
	}

	for _, opt := range options {
		opt(&v)
	}

	c.Lock()
	if c.Data == nil {
		c.Data = make(map[interface{}]*Value)
	}
	c.Data[key] = &v
	c.Unlock()
}

// Get gets a value based on a key.
func (c *Cache) Get(key interface{}) (value interface{}, hit bool) {
	c.RLock()
	valueNode, ok := c.Data[key]
	c.RUnlock()

	if ok {
		value = valueNode.Value
		hit = true
		return
	}
	return
}

// Remove removes a specific key.
func (c *Cache) Remove(key interface{}) (hit bool) {
	c.Lock()
	value, ok := c.Data[key]
	c.Unlock()

	if !ok {
		return
	}

	if value.OnRemove != nil {
		value.OnRemove()
	}
	delete(c.Data, key)
	hit = true
	return
}

// Close stops the cache sweep if it's started.
func (c *Cache) Close() error {
	c.Lock()
	defer c.Unlock()
	if c.SweepCancel == nil {
		return nil
	}

	c.SweepCanceled = make(chan struct{})
	close(c.SweepCancel)
	<-c.SweepCanceled
	return nil
}
