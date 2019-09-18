package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestIsAlpha(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsAlpha('a'))
	assert.True(IsAlpha('z'))
	assert.True(IsAlpha('A'))
	assert.True(IsAlpha('Z'))
	assert.True(IsAlpha('1'))

	assert.False(IsAlpha('-'))
	assert.False(IsAlpha(' '))
	assert.False(IsAlpha('\n'))
	assert.False(IsAlpha('\r'))
	assert.False(IsAlpha('\t'))
}
