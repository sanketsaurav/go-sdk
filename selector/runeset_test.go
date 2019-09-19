package selector

import (
	"testing"

	"github.com/blend/assert/assert"
)

func TestRunSet(t *testing.T) {
	assert := assert.New(t)

	runeset := RuneSet(Dash, Dot, Star)
	assert.NotNil(runeset)
	assert.NotEmpty(runeset)
	assert.Len(runeset, 3)

	_, ok := runeset[Dash]
	assert.True(ok)

	_, ok = runeset[Dot]
	assert.True(ok)

	_, ok = runeset[Star]
	assert.True(ok)

	_, ok = runeset[ForwardSlash]
	assert.False(ok)
}
