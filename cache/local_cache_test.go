package cache

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
)

var (
	_ graceful.Graceful = (*LocalCache)(nil)
)

type itemKey struct{}
type altItemKey struct{}

func TestLocalCache(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache()

	t1 := time.Date(2019, 06, 14, 12, 10, 9, 8, time.UTC)
	c.Set(itemKey{}, "foo", OptValueTimestamp(t1))

	assert.True(c.Has(itemKey{}))
	assert.False(c.Has(altItemKey{}))

	found, ok := c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("foo", found)

	c.Set(altItemKey{}, "alt-bar")
	assert.True(c.Has(altItemKey{}))

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("alt-bar", found)

	t2 := time.Date(2019, 06, 14, 00, 01, 02, 03, time.UTC)
	c.Set(itemKey{}, "bar", OptValueTimestamp(t2))
	assert.Any(c.LRU.Values, func(v interface{}) bool {
		return v.(*Value).Timestamp == t2
	}, "we should have updated the LRU values")

	t3 := time.Date(2019, 06, 14, 12, 01, 02, 03, time.UTC)
	c.Set(altItemKey{}, "alt-bar-2", OptValueTimestamp(t3))
	assert.Any(c.LRU.Values, func(v interface{}) bool {
		return v.(*Value).Timestamp == t3
	}, "we should have updated the LRU values")

	found, ok = c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("bar", found)

	c.Remove(itemKey{})
	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)
	assert.Len(c.LRU.Values, 1)

	c.Set(itemKey{}, "bar", OptValueTimestamp(time.Now().UTC().Add(-time.Hour)))
	assert.True(c.Has(itemKey{}))

	stats := c.Stats()
	assert.Equal(2, stats.Count)
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
		c.Set(nil, "bar")
	}))
	assert.NotNil(try(func() {
		c.Set([]int{}, "bar")
	}))
}

func TestLocalCacheSweep(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache()

	var didSweep, didRemove bool
	c.Set(itemKey{}, "foo",
		OptValueTimestamp(time.Now().UTC().Add(-2*time.Minute)),
		OptValueTTL(time.Minute),
		OptValueOnRemove(func(reason RemovalReason) {
			if reason == ExpiredTTL {
				didSweep = true
			}
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
		OptValueOnRemove(func(reason RemovalReason) {
			if reason == ExpiredTTL {
				close(didSweep)
			}
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

func BenchmarkLocalCache(b *testing.B) {
	for x := 0; x < b.N; x++ {
		benchLocalCache(1024)
	}
}

func benchLocalCache(items int) {
	lc := NewLocalCache()
	for x := 0; x < items; x++ {
		lc.Set(x, strconv.Itoa(x), OptValueTTL(time.Millisecond))
	}
	var value interface{}
	var ok bool
	for x := 0; x < items; x++ {
		value, ok = lc.Get(x)
		if !ok {
			panic("value not found")
		}
		if value.(string) != strconv.Itoa(x) {
			panic("wrong value")
		}
	}
	lc.Sweep(context.Background())
}
