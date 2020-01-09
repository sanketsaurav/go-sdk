package webutil

import (
	"fmt"
	"net"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestIsNetOpError(t *testing.T) {
	assert := assert.New(t)

	assert.False(ex.Filter(IsNetOpError()).Any(fmt.Errorf("foo")))
	assert.True(ex.Filter(IsNetOpError()).Any(&net.OpError{}))
}
