package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

type itemKey struct{}
type altItemKey struct{}

func TestCache(t *testing.T) {
	assert := assert.New(t)

	c := New()
	defer c.Close()

	c.Set(itemKey{}, "foo")

	found, ok := c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("foo", found)

	c.Set(itemKey{}, "bar")

	found, ok = c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("bar", found)

	c.Remove(itemKey{})
	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)
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

func TestCacheKeyPanic(t *testing.T) {
	assert := assert.New(t)

	c := New()
	defer c.Close()

	assert.NotNil(try(func() {
		c.Set("foo", "bar")
	}))
}

func TestCacheSweep(t *testing.T) {
	assert := assert.New(t)

	c := New()
	defer c.Close()

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

	c.Sweep()

	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)
	assert.True(didSweep)
	assert.False(didRemove)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("bar", found)
}

func TestCacheStartSweeping(t *testing.T) {
	assert := assert.New(t)

	c := New(OptSweepInterval(time.Millisecond))
	defer c.Close()

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

	go c.StartSweeping()
	<-c.SweepStarted
	<-didSweep

	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)

}
