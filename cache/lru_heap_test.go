package cache

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestLRUHeap(t *testing.T) {
	assert := assert.New(t)

	t0 := time.Date(2019, 06, 13, 12, 10, 9, 8, time.UTC)
	t1 := time.Date(2019, 06, 14, 12, 10, 9, 8, time.UTC)
	t2 := time.Date(2019, 06, 15, 12, 10, 9, 8, time.UTC)
	t3 := time.Date(2019, 06, 16, 12, 10, 9, 8, time.UTC)
	t4 := time.Date(2019, 06, 17, 12, 10, 9, 8, time.UTC)
	t5 := time.Date(2019, 06, 18, 12, 10, 9, 8, time.UTC)

	h := NewLRUHeap()
	h.Push(&Value{
		Key:       "5",
		Timestamp: t5,
	})
	assert.Len(h.Values, 1)
	h.Push(&Value{
		Key:       "2",
		Timestamp: t2,
	})
	assert.Len(h.Values, 2)
	h.Push(&Value{
		Key:       "3",
		Timestamp: t3,
	})
	assert.Len(h.Values, 3)
	h.Push(&Value{
		Key:       "0",
		Timestamp: t0,
	})
	assert.Len(h.Values, 4)
	h.Push(&Value{
		Key:       "4",
		Timestamp: t4,
	})
	assert.Len(h.Values, 5)
	h.Push(&Value{
		Key:       "1",
		Timestamp: t1,
	})
	assert.Len(h.Values, 6)
	assert.Equal(t0, h.Values[0].Timestamp)

	popped := h.Pop()
	assert.Equal(t0, popped.Timestamp)
	assert.Equal(t1, h.Values[0].Timestamp)

	popped = h.Pop()
	assert.Equal(t1, popped.Timestamp)
	assert.Equal(t2, h.Values[0].Timestamp)

	popped = h.Pop()
	assert.Equal(t2, popped.Timestamp)
	assert.Equal(t3, h.Values[0].Timestamp)

	popped = h.Pop()
	assert.Equal(t3, popped.Timestamp)
	assert.Equal(t4, h.Values[0].Timestamp)

	popped = h.Pop()
	assert.Equal(t4, popped.Timestamp)
	assert.Equal(t5, h.Values[0].Timestamp)

	popped = h.Pop()
	assert.Equal(t5, popped.Timestamp)
	assert.Empty(h.Values)

	popped = h.Pop()
	assert.Nil(popped)
}

func TestLRUHeapConsumeUntil(t *testing.T) {
	assert := assert.New(t)

	t0 := time.Date(2019, 06, 13, 12, 10, 9, 8, time.UTC)
	t1 := time.Date(2019, 06, 14, 12, 10, 9, 8, time.UTC)
	t2 := time.Date(2019, 06, 15, 12, 10, 9, 8, time.UTC)
	t3 := time.Date(2019, 06, 16, 12, 10, 9, 8, time.UTC)
	t4 := time.Date(2019, 06, 17, 12, 10, 9, 8, time.UTC)
	t5 := time.Date(2019, 06, 18, 12, 10, 9, 8, time.UTC)

	h := NewLRUHeap()
	h.Push(&Value{Key: "5", Timestamp: t5})
	h.Push(&Value{Key: "2", Timestamp: t2})
	h.Push(&Value{Key: "3", Timestamp: t3})
	h.Push(&Value{Key: "0", Timestamp: t0})
	h.Push(&Value{Key: "4", Timestamp: t4})
	h.Push(&Value{Key: "1", Timestamp: t1})
	assert.Len(h.Values, 6)

	h.ConsumeUntil(func(v *Value) bool {
		return v.Timestamp.Before(t3)
	})
	assert.Len(h.Values, 3, "consumeUntil should have removed (3) items")
}
