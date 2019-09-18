package selector

import "unicode/utf8"

var (
	_ ValidationRules = (*Kubernetes)(nil)
)

// Kubernetes represent options for validation.
type Kubernetes struct {
	MaxDNSPrefixLen int
	MaxKeyLen       int
	MaxValueLen     int
	NameSymbols     map[rune]struct{}
}

// MaxKeyTotalLen is the maximum total key length
func (k Kubernetes) MaxKeyTotalLen() int {
	return k.MaxDNSPrefixLen + k.MaxKeyLen + 1
}

// IsNameSymbol returns if a given rune is a name symbol.
// Name symbols are non-alphanumeric characters that are allowed
// in either keys or values.
func (k Kubernetes) IsNameSymbol(ch rune) bool {
	if len(k.NameSymbols) > 0 {
		if _, ok := k.NameSymbols[ch]; ok {
			return true
		}
		return false
	}
	switch ch {
	case Dot, Dash, Underscore:
		return true
	}
	return false
}

// CheckKey validates a key.
func (k Kubernetes) CheckKey(key string) (err error) {
	keyLen := len(key)
	if keyLen == 0 {
		err = ErrKeyEmpty
		return
	}
	if keyLen > k.MaxKeyTotalLen() {
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
				err = k.CheckDNS(string(working))
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
	if len(working) > k.MaxKeyLen {
		return ErrKeyTooLong
	}

	return k.CheckName(string(working))
}

// CheckValue returns if the value is valid.
func (k Kubernetes) CheckValue(value string) error {
	if len(value) > k.MaxValueLen {
		return ErrValueTooLong
	}
	return k.CheckName(value)
}

// CheckName returns if a given name is valid.
func (k Kubernetes) CheckName(value string) (err error) {
	valueLen := len(value)
	var state int
	var ch rune
	var width int
	for pos := 0; pos < valueLen; pos += width {
		ch, width = utf8.DecodeRuneInString(value[pos:])
		switch state {
		case 0: //check prefix/suffix
			if !IsAlpha(ch) {
				err = ErrKeyInvalidCharacter
				return
			}
			state = 1
			continue
		case 1:
			if !(k.IsNameSymbol(ch) || ch == BackSlash || IsAlpha(ch)) {
				err = ErrKeyInvalidCharacter
				return
			}
			if pos == valueLen-2 {
				state = 0
			}
			continue
		}
	}
	return
}

// CheckDNS returns if a given string is a valid dns name with a given set of options.
func (k Kubernetes) CheckDNS(value string) (err error) {
	valueLen := len(value)
	if valueLen == 0 {
		err = ErrKeyDNSPrefixEmpty
		return
	}
	if valueLen > k.MaxDNSPrefixLen {
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
			if !IsLowerAlpha(ch) {
				return ErrKeyInvalidCharacter
			}
			state = 1
			continue
		case 1:
			if ch == Underscore {
				err = ErrKeyInvalidCharacter
				return
			}
			if k.IsNameSymbol(ch) {
				state = 2
				continue
			}
			if !IsLowerAlpha(ch) {
				err = ErrKeyInvalidCharacter
				return
			}
			if pos == valueLen-2 {
				state = 0
			}
			continue
		case 2: // we've hit a dot, dash, or underscore that can't repeat
			if !IsLowerAlpha(ch) {
				err = ErrKeyInvalidCharacter
				return
			}
			if pos == valueLen-2 {
				state = 0
			}
			state = 1
		}
	}
	return nil
}
