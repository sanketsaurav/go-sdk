package ex

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFilters(t *testing.T) {
	assert := assert.New(t)

	isFoo := FilterIs(fmt.Errorf("foo"))
	isBar := FilterIs(fmt.Errorf("bar"))

	filters := Filter(isFoo, isBar)
	assert.True(filters.Any(fmt.Errorf("foo")))
	assert.True(filters.Any(fmt.Errorf("bar")))

	assert.False(filters.All(fmt.Errorf("foo")))
	assert.False(filters.All(fmt.Errorf("bar")))
}
