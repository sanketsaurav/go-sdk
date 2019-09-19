package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewRuneSet(t *testing.T) {
	assert := assert.New(t)

	runeset := NewRuneSet(DefaultNameSymbols...)
	assert.NotNil(runeset)
	assert.NotEmpty(runeset)
	assert.Len(runeset, 3)

	assert.True(runeset.Has(Dash))
	assert.True(runeset.Has(Dot))
	assert.True(runeset.Has(Underscore))
	assert.False(runeset.Has(ForwardSlash))
}
