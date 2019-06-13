package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMemoize(t *testing.T) {
	assert := assert.New(t)

	var calls int
	identity := func(args interface{}) (interface{}, error) {
		calls++
		return args, nil
	}

	mIdentity := Memoize(identity)

	res, err := mIdentity.Call("foo")
	assert.Nil(err)
	assert.Equal("foo", res)
	assert.Equal(1, calls)

	res, err = mIdentity.Call("bar")
	assert.Nil(err)
	assert.Equal("bar", res)
	assert.Equal(2, calls)

	res, err = mIdentity.Call("foo")
	assert.Nil(err)
	assert.Equal("foo", res)
	assert.Equal(2, calls)
}

func TestMemoizeTTL(t *testing.T) {
	assert := assert.New(t)

	var calls int
	identity := func(args interface{}) (interface{}, error) {
		calls++
		return args, nil
	}

	mIdentity := Memoize(identity, OptMemoizeTTL(time.Millisecond))
	assert.NotZero(mIdentity.TTL)

	res, err := mIdentity.Call("foo")
	assert.Nil(err)
	assert.Equal("foo", res)
	assert.Equal(1, calls)

	res, err = mIdentity.Call("foo")
	assert.Nil(err)
	assert.Equal("foo", res)
	assert.Equal(1, calls)

	<-time.After(2 * time.Millisecond)

	res, err = mIdentity.Call("foo")
	assert.Nil(err)
	assert.Equal("foo", res)
	assert.Equal(2, calls)
}

func TestMemoizeErr(t *testing.T) {
	assert := assert.New(t)

	var calls int
	identity := func(args interface{}) (interface{}, error) {
		calls++
		return args, fmt.Errorf("this is only a test")
	}

	mIdentity := Memoize(identity)

	res, err := mIdentity.Call("foo")
	assert.NotNil(err)
	assert.Equal("foo", res)
	assert.Equal(1, calls)

	res, err = mIdentity.Call("bar")
	assert.NotNil(err)
	assert.Equal("bar", res)
	assert.Equal(2, calls)

	res, err = mIdentity.Call("foo")
	assert.NotNil(err)
	assert.Equal("foo", res)
	assert.Equal(2, calls)
}

func TestMemoizeComplex(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		ID   int
		Name string
	}

	var calls int
	identity := func(args interface{}) (interface{}, error) {
		calls++
		return args, nil
	}

	mIdentity := Memoize(identity)

	res, err := mIdentity.Call(args{ID: 1, Name: "one"})
	assert.Nil(err)
	assert.Equal("one", res.(args).Name)
	assert.Equal(1, calls)

	res, err = mIdentity.Call(args{ID: 2, Name: "two"})
	assert.Nil(err)
	assert.Equal("two", res.(args).Name)
	assert.Equal(2, calls)

	res, err = mIdentity.Call(args{ID: 1, Name: "one"})
	assert.Nil(err)
	assert.Equal("one", res.(args).Name)
	assert.Equal(2, calls)
}
