package selector

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/blend/go-sdk/ex"
)

// Parser parses a selector incrementally.
type Parser struct {
	// SkipValidation skips all validation.
	SkipValidation bool
	// AllowSymbolPrefix allows key and values to start with
	// a rune from the name symbols.
	AllowSymbolPrefix bool
	// AllowSymbolSuffix allows keys and values to end with
	// a rune from the name symbols.
	AllowSymbolSuffix bool
	// AllowSymbolRepeats allows keys and values to
	// include repeated symbols.
	AllowSymbolRepeats bool
	// MaxDNSPrefixLen is the maximum dns prefix length.
	MaxDNSPrefixLen int
	// MaxKeyLen is the maximum key length.
	MaxKeyLen int
	// MaxValueLen is the maximum value length.
	MaxValueLen int
	// NameSymbols are the symbols that can be in name
	// (either keys or values) in addition to alphanumeric characters.
	NameSymbols map[rune]struct{}

	// s stores the string to be tokenized
	s string
	// pos is the position currently tokenized
	pos int
	// m is an optional mark
	m int
}

// Parse does the actual parsing.
func (p *Parser) Parse() (Selector, error) {
	p.s = strings.TrimSpace(p.s)
	if len(p.s) == 0 {
		return Any{}, nil
	}

	var b rune
	var selector Selector
	var err error
	var op string

	// loop over "clauses"
	// clauses are separated by commas and grouped logically as "ands"
	for {
		// sniff the !haskey form
		b = p.current()
		if b == Bang {
			p.advance() // we aren't going to use the '!'
			selector = p.addAnd(selector, p.notHasKey(p.readWord()))
			if p.done() {
				break
			}
			continue
		}

		// we're done peeking the first char
		// read the next word as a key
		// and check it is valid.
		key := p.readWord()
		if err = p.CheckKey(key); err != nil {
			return nil, err
		}

		p.mark()

		// check if the next character after the word is a comma
		// this indicates it's a "key" form, or existence check on a key
		b = p.skipToComma()
		if b == Comma || p.isTerminator(b) || p.done() {
			selector = p.addAnd(selector, p.hasKey(key))
			p.advance()
			if p.done() {
				break
			}
			continue
		} else {
			p.popMark()
		}

		op, err = p.readOp()
		if err != nil {
			return nil, err
		}

		var subSelector Selector
		switch op {
		case OpEquals, OpDoubleEquals:
			subSelector, err = p.equals(key)
			if err != nil {
				return nil, err
			}
			selector = p.addAnd(selector, subSelector)
		case OpNotEquals:
			subSelector, err = p.notEquals(key)
			if err != nil {
				return nil, err
			}
			selector = p.addAnd(selector, subSelector)
		case OpIn:
			subSelector, err = p.in(key)
			if err != nil {
				return nil, err
			}
			selector = p.addAnd(selector, subSelector)
		case OpNotIn:
			subSelector, err = p.notIn(key)
			if err != nil {
				return nil, err
			}
			selector = p.addAnd(selector, subSelector)
		default:
			return nil, ErrInvalidOperator
		}

		b = p.skipToComma()
		if b == Comma {
			p.advance()
			if p.done() {
				break
			}
			continue
		}

		// these two are effectively the same
		if p.isTerminator(b) || p.done() {
			break
		}

		return nil, ErrInvalidSelector
	}

	return selector, nil
}

// MaxDNSPrefixLenOrDefault returns the max dns prefix length or a default.
func (p *Parser) MaxDNSPrefixLenOrDefault() int {
	if p.MaxDNSPrefixLen > 0 {
		return p.MaxDNSPrefixLen
	}
	return DefaultMaxDNSPrefixLen
}

// MaxKeyLenOrDefault returns the max key length or a default.
func (p *Parser) MaxKeyLenOrDefault() int {
	if p.MaxKeyLen > 0 {
		return p.MaxKeyLen
	}
	return DefaultMaxKeyLen
}

// MaxValueLenOrDefault returns the max value length or a default.
func (p *Parser) MaxValueLenOrDefault() int {
	if p.MaxValueLen > 0 {
		return p.MaxValueLen
	}
	return DefaultMaxValueLen
}

// MaxKeyTotalLen is the maximum total key length
func (p *Parser) MaxKeyTotalLen() int {
	return p.MaxDNSPrefixLenOrDefault() + p.MaxKeyLenOrDefault() + 1
}

// CheckKey validates a key.
func (p *Parser) CheckKey(key string) (err error) {
	if p.SkipValidation {
		return nil
	}
	keyLen := len(key)
	if keyLen == 0 {
		err = ErrKeyEmpty
		return
	}
	if keyLen > p.MaxKeyTotalLen() {
		err = ErrKeyTooLong
		return
	}

	var working []rune
	var state int
	var ch rune
	var width int
	for pos := 0; pos < keyLen; pos += width {
		ch, width = utf8.DecodeRuneInString(key[pos:])
		switch state {
		case 0: // collect dns prefix or key
			if ch == ForwardSlash {
				err = p.CheckDNS(string(working))
				if err != nil {
					return
				}
				working = nil
				state = 1
				continue
			}
		}
		working = append(working, ch)
		continue
	}

	if len(working) == 0 {
		return ErrKeyEmpty
	}
	if len(working) > p.MaxKeyLenOrDefault() {
		return ErrKeyTooLong
	}

	return p.CheckName(string(working))
}

// CheckValue returns if the value is valid.
func (p *Parser) CheckValue(value string) error {
	if p.SkipValidation {
		return nil
	}
	if len(value) > p.MaxValueLenOrDefault() {
		return ex.New(ErrValueTooLong, ex.OptMessagef("must be less than %d", p.MaxValueLenOrDefault()))
	}
	return p.CheckName(value)
}

// CheckName returns if a given name is valid.
func (p *Parser) CheckName(value string) (err error) {
	valueLen := len(value)
	var state int
	var ch rune
	var width int
	for pos := 0; pos < valueLen; pos += width {
		ch, width = utf8.DecodeRuneInString(value[pos:])
		switch state {
		case 0: //check prefix
			if !p.isAlpha(ch) && !p.AllowSymbolPrefix {
				err = p.errKeyInvalidCharacter()
				return
			}
			state = 1
			continue
		case 1:
			if !(p.isNameSymbol(ch) || ch == BackSlash || p.isAlpha(ch)) {
				err = p.errKeyInvalidCharacter()
				return
			}
			if pos == valueLen-2 {
				state = 2
			}
			continue
		case 2: //check suffix
			if !p.isAlpha(ch) && !p.AllowSymbolSuffix {
				err = p.errKeyInvalidCharacter()
				return
			}
			state = 1
			continue
		}
	}
	return
}

// CheckDNS returns if a given string is a valid dns name with a given set of options.
func (p *Parser) CheckDNS(value string) (err error) {
	valueLen := len(value)
	if valueLen == 0 {
		err = ErrKeyDNSPrefixEmpty
		return
	}
	if valueLen > p.MaxDNSPrefixLenOrDefault() {
		err = ErrKeyDNSPrefixTooLong
		return
	}
	var state int
	var ch rune
	var width int
	for pos := 0; pos < valueLen; pos += width {
		ch, width = utf8.DecodeRuneInString(value[pos:])
		switch state {
		case 0: //check prefix | suffix
			if !p.isLowerAlpha(ch) && !p.AllowSymbolPrefix {
				return p.errKeyInvalidCharacter()
			}
			state = 1
			continue
		case 1:
			if ch == Underscore {
				err = p.errKeyInvalidCharacter()
				return
			}
			if p.isNameSymbol(ch) && !p.AllowSymbolRepeats {
				state = 2
				continue
			}
			if !p.isLowerAlpha(ch) {
				err = p.errKeyInvalidCharacter()
				return
			}
			if pos == valueLen-2 {
				state = 3
			}
			continue
		case 2: // we've hit a dot, dash, or underscore that can't repeat
			if !p.isLowerAlpha(ch) {
				err = p.errKeyInvalidCharacter()
				return
			}
			if pos == valueLen-2 {
				state = 3
			}
			state = 1
		case 3: //check suffix
			if !p.isLowerAlpha(ch) && !p.AllowSymbolSuffix {
				return p.errKeyInvalidCharacter()
			}
			state = 1
			continue
		}
	}
	return nil
}

// done indicates the cursor is past the usable length of the string.
func (p *Parser) done() bool {
	return p.pos == len(p.s)
}

// mark sets a mark at the current position.
func (p *Parser) mark() {
	p.m = p.pos
}

// popMark moves the cursor back to the previous mark.
func (p *Parser) popMark() {
	if p.m > 0 {
		p.pos = p.m
	}
	p.m = 0
}

// read returns the rune currently lexed, and advances the position.
func (p *Parser) read() (r rune) {
	var width int
	if p.pos < len(p.s) {
		r, width = utf8.DecodeRuneInString(p.s[p.pos:])
		p.pos += width
	}
	return r
}

// current returns the rune at the current position.
func (p *Parser) current() (r rune) {
	r, _ = utf8.DecodeRuneInString(p.s[p.pos:])
	return
}

// advance moves the cursor forward one rune.
func (p *Parser) advance() {
	if p.pos < len(p.s) {
		_, width := utf8.DecodeRuneInString(p.s[p.pos:])
		p.pos += width
	}
}

// prev moves the cursor back a rune.
func (p *Parser) prev() {
	if p.pos > 0 {
		p.pos--
	}
}

// readOp reads a valid operator.
// valid operators include:
// [ =, ==, !=, in, notin ]
// errors if it doesn't read one of the above, or there is another structural issue.
func (p *Parser) readOp() (string, error) {
	// skip preceding whitespace
	p.skipWhiteSpace()

	var state int
	var ch rune
	var op []rune
	for {
		ch = p.current()

		switch state {
		case 0: // initial state, determine what op we're reading for
			if ch == Equal {
				state = 1
				break
			}
			if ch == Bang {
				state = 2
				break
			}
			if ch == 'i' {
				state = 6
				break
			}
			if ch == 'n' {
				state = 7
				break
			}
			return "", ErrInvalidOperator
		case 1: // =
			if unicode.IsSpace(ch) || p.isAlpha(ch) || p.isNameSymbol(ch) || ch == Comma {
				return string(op), nil
			}
			if ch == Equal {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 2: // !
			if ch == Equal {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 6: // in
			if ch == 'n' {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 7: // o
			if ch == 'o' {
				state = 8
				break
			}
			return "", ErrInvalidOperator
		case 8: // t
			if ch == 't' {
				state = 9
				break
			}
			return "", ErrInvalidOperator
		case 9: // i
			if ch == 'i' {
				state = 10
				break
			}
			return "", ErrInvalidOperator
		case 10: // n
			if ch == 'n' {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		}

		op = append(op, ch)
		p.advance()

		if p.done() {
			return string(op), nil
		}
	}
}

// readWord skips whitespace, then reads a word until whitespace or a token.
// it will leave the cursor on the next char after the word, i.e. the space or token.
func (p *Parser) readWord() string {
	// skip preceding whitespace
	p.skipWhiteSpace()

	var word []rune
	var ch rune
	for {
		ch = p.current()

		if unicode.IsSpace(ch) {
			return string(word)
		}
		if p.isSelectorSymbol(ch) {
			return string(word)
		}

		word = append(word, ch)
		p.advance()

		if p.done() {
			return string(word)
		}
	}
}

func (p *Parser) readCSV() (results []string, err error) {
	var word []rune
	var ch rune
	var state int

	// skip preceding whitespace
	p.skipWhiteSpace()

	for {
		ch = p.current()

		if p.done() {
			err = ErrInvalidSelector
			return
		}

		switch state {
		case 0: // leading paren
			if ch == OpenParens {
				state = 2 // spaces or alphas
				p.advance()
				continue
			}
			// not open parens, bail
			err = ErrInvalidSelector
			return
		case 1: // alphas (in word)

			if ch == Comma {
				if len(word) > 0 {
					if err = p.CheckValue(string(word)); err != nil {
						return nil, err
					}
					results = append(results, string(word))
					word = nil
				}
				state = 2 // from comma
				p.advance()
				continue
			}

			if ch == CloseParens {
				if len(word) > 0 {
					if err = p.CheckValue(string(word)); err != nil {
						return nil, err
					}
					results = append(results, string(word))
				}
				p.advance()
				return
			}

			if unicode.IsSpace(ch) {
				state = 3
				p.advance()
				continue
			}

			word = append(word, ch)
			p.advance()
			continue

		case 2: //whitespace after symbol

			if ch == CloseParens {
				p.advance()
				return
			}

			if unicode.IsSpace(ch) {
				p.advance()
				continue
			}

			if ch == Comma {
				p.advance()
				continue
			}

			if p.isAlpha(ch) {
				state = 1
				continue
			}

			err = ErrInvalidSelector
			return

		case 3: //whitespace after alpha

			if ch == CloseParens {
				if len(word) > 0 {
					if err = p.CheckValue(string(word)); err != nil {
						return nil, err
					}
					results = append(results, string(word))
				}
				p.advance()
				return
			}

			if unicode.IsSpace(ch) {
				p.advance()
				continue
			}

			if ch == Comma {
				if len(word) > 0 {
					if err = p.CheckValue(string(word)); err != nil {
						return nil, err
					}
					results = append(results, string(word))
					word = nil
				}
				p.advance()
				state = 2
				continue
			}

			err = ErrInvalidSelector
			return
		}
	}
}

func (p *Parser) skipWhiteSpace() {
	if p.done() {
		return
	}
	var ch rune
	for {
		ch = p.current()
		if !unicode.IsSpace(ch) {
			return
		}
		p.advance()
		if p.done() {
			return
		}
	}
}

func (p *Parser) skipToComma() (ch rune) {
	if p.done() {
		return
	}
	for {
		ch = p.current()
		if ch == Comma {
			return
		}
		if !unicode.IsSpace(ch) {
			return
		}
		p.advance()
		if p.done() {
			return
		}
	}
}

func (p *Parser) errKeyInvalidCharacter() error {
	var nameSymbols []string
	if len(p.NameSymbols) > 0 {
		for ch := range p.NameSymbols {
			if ch == Dot || ch == Star {
				nameSymbols = append(nameSymbols, "\\"+string(ch))
			} else {
				nameSymbols = append(nameSymbols, string(ch))
			}
		}
	} else {
		for _, ch := range DefaultNameSymbols {
			if ch == Dot || ch == Star {
				nameSymbols = append(nameSymbols, "\\"+string(ch))
			} else {
				nameSymbols = append(nameSymbols, string(ch))
			}
		}
	}
	sort.Strings(nameSymbols)
	msg := fmt.Sprintf("regex used: [A-Za-z0-9%s]", strings.Join(nameSymbols, ""))
	return &ex.Ex{
		Class:   ErrKeyInvalidCharacter,
		Message: msg,
	}
}

func (p *Parser) isNameSymbol(ch rune) bool {
	if len(p.NameSymbols) > 0 {
		if _, ok := p.NameSymbols[ch]; ok {
			return true
		}
		return false
	}
	// handle defaults
	switch ch {
	case Dot, Dash, Underscore:
		return true
	}
	return false
}

// isSelectorSymbol returns if a rune is allowed in a
// selector phrase [i.e. ! =  (  ) , ]
func (p *Parser) isSelectorSymbol(ch rune) bool {
	switch ch {
	case Equal, Bang, OpenParens, CloseParens, Comma:
		return true
	}
	return false
}

// IsSymbol returns if a rune is a symbol.
func (p *Parser) isSymbol(ch rune) bool {
	return (int(ch) >= int(Bang) && int(ch) <= int(ForwardSlash)) ||
		(int(ch) >= int(Colon) && int(ch) <= int(At)) ||
		(int(ch) >= int(OpenBracket) && int(ch) <= int(BackTick)) ||
		(int(ch) >= int(OpenCurly) && int(ch) <= int(Tilde))
}

// isLowerAlpha returns if a rune is a lower cased alphanumeric.
func (p *Parser) isLowerAlpha(ch rune) bool {
	if unicode.IsLetter(ch) {
		return unicode.IsLower(ch)
	}
	return p.isAlpha(ch)
}

// isAlpha returns if a rune is not a space, control or symbol.
func (p *Parser) isAlpha(ch rune) bool {
	return !unicode.IsSpace(ch) && !unicode.IsControl(ch) && !p.isSymbol(ch)
}

func (p *Parser) isTerminator(ch rune) bool {
	return ch == 0
}

//
// selector node constructors
//

// addAnd starts grouping selectors into a high level `and`, returning the aggregate selector.
func (p *Parser) addAnd(current, next Selector) Selector {
	if current == nil {
		return next
	}
	if typed, isTyped := current.(And); isTyped {
		return append(typed, next)
	}
	return And([]Selector{current, next})
}

func (p *Parser) hasKey(key string) Selector {
	return HasKey(key)
}

func (p *Parser) notHasKey(key string) Selector {
	return NotHasKey(key)
}

func (p *Parser) equals(key string) (Selector, error) {
	value := p.readWord()
	if err := p.CheckValue(value); err != nil {
		return nil, err
	}
	return Equals{Key: key, Value: value}, nil
}

func (p *Parser) notEquals(key string) (Selector, error) {
	value := p.readWord()
	if err := p.CheckValue(value); err != nil {
		return nil, err
	}
	return NotEquals{Key: key, Value: value}, nil
}

func (p *Parser) in(key string) (Selector, error) {
	csv, err := p.readCSV()
	if err != nil {
		return nil, err
	}
	return In{Key: key, Values: csv}, nil
}

func (p *Parser) notIn(key string) (Selector, error) {
	csv, err := p.readCSV()
	if err != nil {
		return nil, err
	}
	return NotIn{Key: key, Values: csv}, nil
}
