package cache

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func tv(index int) *Value {
	return &Value{
		Key:       index,
		Timestamp: time.Date(2019, 06, 14, 12, index, 0, 0, time.UTC),
	}
}

func TestLRUQueue(t *testing.T) {
	assert := assert.New(t)

	buffer := NewLRUQueue()

	buffer.Push(tv(1))
	assert.Equal(1, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(1, buffer.PeekBack().Key)

	buffer.Push(tv(2))
	assert.Equal(2, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(2, buffer.PeekBack().Key)

	buffer.Push(tv(3))
	assert.Equal(3, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(3, buffer.PeekBack().Key)

	buffer.Push(tv(4))
	assert.Equal(4, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(4, buffer.PeekBack().Key)

	buffer.Push(tv(5))
	assert.Equal(5, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(5, buffer.PeekBack().Key)

	buffer.Push(tv(6))
	assert.Equal(6, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(6, buffer.PeekBack().Key)

	buffer.Push(tv(7))
	assert.Equal(7, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(7, buffer.PeekBack().Key)

	buffer.Push(tv(8))
	assert.Equal(8, buffer.Len())
	assert.Equal(1, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value := buffer.Pop()
	assert.Equal(1, value.Key)
	assert.Equal(7, buffer.Len())
	assert.Equal(2, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(2, value.Key)
	assert.Equal(6, buffer.Len())
	assert.Equal(3, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(3, value.Key)
	assert.Equal(5, buffer.Len())
	assert.Equal(4, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(4, value.Key)
	assert.Equal(4, buffer.Len())
	assert.Equal(5, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(5, value.Key)
	assert.Equal(3, buffer.Len())
	assert.Equal(6, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(6, value.Key)
	assert.Equal(2, buffer.Len())
	assert.Equal(7, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(7, value.Key)
	assert.Equal(1, buffer.Len())
	assert.Equal(8, buffer.Peek().Key)
	assert.Equal(8, buffer.PeekBack().Key)

	value = buffer.Pop()
	assert.Equal(8, value.Key)
	assert.Equal(0, buffer.Len())
	assert.Nil(buffer.Peek())
	assert.Nil(buffer.PeekBack())
}

func TestLRUQueueClear(t *testing.T) {
	assert := assert.New(t)

	buffer := NewLRUQueue()
	for x := 0; x < 8; x++ {
		buffer.Push(tv(x))
	}
	assert.Equal(8, buffer.Len())
	buffer.Clear()
	assert.Equal(0, buffer.Len())
	assert.Nil(buffer.Peek())
	assert.Nil(buffer.PeekBack())
}

func TestLRUQueueConsumeUntil(t *testing.T) {
	assert := assert.New(t)

	buffer := NewLRUQueue()

	for x := 1; x < 17; x++ {
		buffer.Push(tv(x))
	}

	called := 0
	buffer.ConsumeUntil(func(v *Value) bool {
		if v.Key.(int) > 10 {
			return false
		}
		if v.Key.(int) == (called + 1) {
			called++
		}
		return true
	})

	assert.Equal(10, called)
}
