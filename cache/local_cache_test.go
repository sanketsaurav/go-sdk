package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

type itemKey struct{}
type altItemKey struct{}

func TestLocalCache(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache()

	c.Set(itemKey{}, "foo")

	assert.True(c.Has(itemKey{}))
	assert.False(c.Has(altItemKey{}))

	found, ok := c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("foo", found)

	c.Set(itemKey{}, "bar")
	assert.True(c.Has(itemKey{}))

	found, ok = c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("bar", found)

	c.Remove(itemKey{})
	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)

	c.Set(itemKey{}, "bar", OptValueTimestamp(time.Now().UTC().Add(-time.Hour)))
	assert.True(c.Has(itemKey{}))

	stats := c.Stats()
	assert.Equal(1, stats.Count)
	assert.NotZero(stats.MaxAge)
	assert.NotZero(stats.SizeBytes)
}

func try(action func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	action()
	return
}

func TestLocalCacheKeyPanic(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache()

	assert.NotNil(try(func() {
		c.Set("foo", "bar")
	}))
}

func TestLocalCacheSweep(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache()

	var didSweep, didRemove bool
	c.Set(itemKey{}, "foo",
		OptValueTimestamp(time.Now().UTC().Add(-2*time.Minute)),
		OptValueTTL(time.Minute),
		OptValueOnSweep(func() {
			didSweep = true
		}),
		OptValueOnRemove(func() {
			didRemove = true
		}),
	)
	found, ok := c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("foo", found)

	c.Set(altItemKey{}, "bar",
		OptValueTTL(time.Minute),
	)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("bar", found)

	c.Sweep(context.Background())

	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)
	assert.True(didSweep)
	assert.False(didRemove)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("bar", found)
}

func TestLocalCacheStartSweeping(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache(OptLocalCacheSweepInterval(time.Millisecond))

	didSweep := make(chan struct{})
	c.Set(itemKey{}, "a value",
		OptValueTTL(time.Microsecond),
		OptValueOnSweep(func() {
			close(didSweep)
		}),
	)

	found, ok := c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("a value", found)

	c.Set(altItemKey{}, "bar",
		OptValueTTL(time.Minute),
	)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("bar", found)

	go c.Start()
	<-c.NotifyStarted()
	defer c.Stop()
	<-didSweep

	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("bar", found)
}
