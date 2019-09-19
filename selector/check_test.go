package selector

import (
	"strings"
	"testing"

	assert "github.com/blend/go-sdk/assert"
)

func TestCheckLabels(t *testing.T) {
	assert := assert.New(t)

	goodLabels := Labels{"foo": "bar", "foo.com/bar": "baz"}
	assert.Nil(CheckLabels(goodLabels))
	badLabels := Labels{"foo": "bar", "_foo.com/bar": "baz"}
	assert.NotNil(CheckLabels(badLabels))
}

func TestCheckKeyK8S(t *testing.T) {
	assert := assert.New(t)

	p := new(Parser)

	values := []string{
		// the "good" cases
		"simple",
		"now-with-dashes",
		"1-starts-with-num",
		"1234",
		"simple/simple",
		"now-with-dashes/simple",
		"now-with-dashes/now-with-dashes",
		"now.with.dots/simple",
		"now-with.dashes-and.dots/simple",
		"1-num.2-num/3-num",
		"1234/5678",
		"1.2.3.4/5678",
		"Uppercase_Is_OK_123",
		"example.com/Uppercase_Is_OK_123",
		"requests.storage-foo",
		strings.Repeat("a", 63),
		strings.Repeat("a", 253) + "/" + strings.Repeat("b", 63),
	}
	badValues := []string{
		// the "bad" cases
		"nospecialchars%^=@",
		"cantendwithadash-",
		"-cantstartwithadash-",
		"only/one/slash",
		"Example.com/abc",
		"example_com/abc",
		"example.com/",
		"/simple",
		strings.Repeat("a", 64),
		strings.Repeat("a", 254) + "/abc",
	}
	for _, val := range values {
		assert.Nil(p.CheckKey(val))
	}
	for _, val := range badValues {
		assert.NotNil(p.CheckKey(val), val)
	}
}

func TestCheckValue(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckValue(""), "should not error on empty values")
	assert.Nil(CheckValue("foo"))
	assert.Nil(CheckValue("bar_baz"))
	assert.NotNil(CheckValue("_bar_baz"))
	assert.NotNil(CheckValue("bar_baz_"))
	assert.NotNil(CheckValue("_bar_baz_"))
}
