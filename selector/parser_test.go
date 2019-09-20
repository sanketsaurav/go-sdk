package selector

import (
	"fmt"
	"strings"
	"testing"

	assert "github.com/blend/go-sdk/assert"
)

func TestParserIsAlpha(t *testing.T) {
	assert := assert.New(t)

	p := new(Parser)
	assert.True(p.isAlpha('a'))
	assert.True(p.isAlpha('z'))
	assert.True(p.isAlpha('A'))
	assert.True(p.isAlpha('Z'))
	assert.True(p.isAlpha('1'))

	assert.False(p.isAlpha('-'))
	assert.False(p.isAlpha(' '))
	assert.False(p.isAlpha('\n'))
	assert.False(p.isAlpha('\r'))
	assert.False(p.isAlpha('\t'))
}

func TestParserSkipWhitespace(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "foo    != bar    ", pos: 3}
	assert.Equal(" ", string(l.current()))
	l.skipWhiteSpace()
	assert.Equal(7, l.pos)
	assert.Equal("!", string(l.current()))
	l.pos = 14
	assert.Equal(" ", string(l.current()))
	l.skipWhiteSpace()
	assert.Equal(len(l.s), l.pos)
}

func TestParserReadWord(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "foo != bar"}
	assert.Equal("foo", l.readWord())
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "foo,"}
	assert.Equal("foo", l.readWord())
	assert.Equal(",", string(l.current()))

	l = &Parser{s: "foo,bar,baz", pos: 4}
	assert.Equal("b", string(l.current()))
	assert.Equal("bar", l.readWord())
	assert.Equal(",", string(l.current()))

	l = &Parser{s: "foo"}
	assert.Equal("foo", l.readWord())
	assert.True(l.done())
}

func TestParserReadOp(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "!= bar"}
	op, err := l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "!=bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.Equal("b", string(l.current()))

	l = &Parser{s: "!=bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.Equal("b", string(l.current()))

	l = &Parser{s: "!="}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.True(l.done())

	l = &Parser{s: "= bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("=", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "=bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("=", op)
	assert.Equal("b", string(l.current()))

	l = &Parser{s: "== bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("==", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "==bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("==", op)
	assert.Equal("b", string(l.current()))

	l = &Parser{s: "in (foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("in", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "in(foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("in", op)
	assert.Equal("(", string(l.current()))

	l = &Parser{s: "notin (foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("notin", op)
	assert.Equal(" ", string(l.current()))

	l = &Parser{s: "notin(foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("notin", op)
	assert.Equal("(", string(l.current()))
}

func TestParserReadCSV(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "(bar, baz, biz)"}
	words, err := l.readCSV()
	assert.Nil(err)
	assert.Len(words, 3, strings.Join(words, ","))
	assert.Equal("bar", words[0])
	assert.Equal("baz", words[1])
	assert.Equal("biz", words[2])
	assert.True(l.done())

	l = &Parser{s: "(bar,baz,biz)"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Len(words, 3, strings.Join(words, ","))
	assert.Equal("bar", words[0])
	assert.Equal("baz", words[1])
	assert.Equal("biz", words[2])
	assert.True(l.done())

	l = &Parser{s: "(bar, buzz, baz"}
	words, err = l.readCSV()
	assert.NotNil(err)

	l = &Parser{s: "()"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Empty(words)
	assert.True(l.done())

	l = &Parser{s: "(), thing=after"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Empty(words)
	assert.Equal(",", string(l.current()))

	l = &Parser{s: "(foo, bar), buzz=light"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Len(words, 2)
	assert.Equal("foo", words[0])
	assert.Equal("bar", words[1])
	assert.Equal(",", string(l.current()))

	l = &Parser{s: "(test, space are bad)"}
	words, err = l.readCSV()
	assert.NotNil(err)
}

func TestParserHasKey(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "foo"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(HasKey)
	assert.True(isTyped)
	assert.Equal("foo", string(typed))
}

func TestParserNotHasKey(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "!foo"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(NotHasKey)
	assert.True(isTyped)
	assert.Equal("foo", string(typed))
}

func TestParserEquals(t *testing.T) {
	assert := assert.New(t)

	l := &Parser{s: "foo = bar"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(Equals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)

	l = &Parser{s: "foo=bar"}
	valid, err = l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped = valid.(Equals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)
}

func TestParserDoubleEquals(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "foo == bar"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(Equals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)
}

func TestParserNotEquals(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "foo != bar"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(NotEquals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)
}

func TestParserIn(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: "foo in (bar, baz)"}
	valid, err := l.Parse()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(In)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Len(typed.Values, 2)
	assert.Equal("bar", typed.Values[0])
	assert.Equal("baz", typed.Values[1])
}

func TestParserLex(t *testing.T) {
	assert := assert.New(t)
	l := &Parser{s: ""}
	_, err := l.Parse()
	assert.Nil(err)
}

func TestParserCheckKey(t *testing.T) {
	assert := assert.New(t)

	p := new(Parser)
	assert.Nil(p.CheckKey("foo"))
	assert.Nil(p.CheckKey("bar/foo"))
	assert.Nil(p.CheckKey("bar.io/foo"))
	assert.NotNil(p.CheckKey("_foo"))
	assert.NotNil(p.CheckKey("-foo"))
	assert.NotNil(p.CheckKey("foo-"))
	assert.NotNil(p.CheckKey("foo_"))
	assert.NotNil(p.CheckKey("bar/foo/baz"))

	assert.NotNil(p.CheckKey(""), "should error on empty keys")

	assert.NotNil(p.CheckKey("/foo"), "should error on empty dns prefixes")
	superLongDNSPrefixed := fmt.Sprintf("%s/%s", strings.Repeat("a", p.MaxDNSPrefixLenOrDefault()), strings.Repeat("a", p.MaxKeyLenOrDefault()))
	assert.Nil(p.CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", p.MaxDNSPrefixLenOrDefault()+1), strings.Repeat("a", p.MaxKeyLenOrDefault()))
	assert.NotNil(p.CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", p.MaxDNSPrefixLenOrDefault()+1), strings.Repeat("a", p.MaxKeyLenOrDefault()+1))
	assert.NotNil(p.CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", p.MaxDNSPrefixLenOrDefault()), strings.Repeat("a", p.MaxKeyLenOrDefault()+1))
	assert.NotNil(p.CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
}

func TestParserCheckKeyK8S(t *testing.T) {
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

func TestParserCheckValue(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckValue(""), "should not error on empty values")
	assert.Nil(CheckValue("foo"))
	assert.Nil(CheckValue("bar_baz"))
	assert.NotNil(CheckValue("_bar_baz"))
	assert.NotNil(CheckValue("bar_baz_"))
	assert.NotNil(CheckValue("_bar_baz_"))
}

func TestParserCheckLabels(t *testing.T) {
	assert := assert.New(t)

	goodLabels := Labels{"foo": "bar", "foo.com/bar": "baz"}
	assert.Nil(CheckLabels(goodLabels))
	badLabels := Labels{"foo": "bar", "_foo.com/bar": "baz"}
	assert.NotNil(CheckLabels(badLabels))
}
