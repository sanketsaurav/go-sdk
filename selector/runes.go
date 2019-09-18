package selector

import "unicode"

// IsSelectorSymbol returns if a rune is allowed in a
// selector phrase [i.e. ! =  (  ) , ]
func IsSelectorSymbol(ch rune) bool {
	switch ch {
	case Equal, Bang, OpenParens, CloseParens, Comma:
		return true
	}
	return false
}

// IsSymbol returns if a rune is a symbol.
func IsSymbol(ch rune) bool {
	return (int(ch) >= int(Bang) && int(ch) <= int(ForwardSlash)) ||
		(int(ch) >= int(Colon) && int(ch) <= int(At)) ||
		(int(ch) >= int(OpenBracket) && int(ch) <= int(BackTick)) ||
		(int(ch) >= int(OpenCurly) && int(ch) <= int(Tilde))
}

// IsLowerAlpha returns if a rune is a lower cased alphanumeric.
func IsLowerAlpha(ch rune) bool {
	if unicode.IsLetter(ch) {
		return unicode.IsLower(ch)
	}
	return IsAlpha(ch)
}

// IsAlpha returns if a rune is not a space, control or symbol.
func IsAlpha(ch rune) bool {
	return !unicode.IsSpace(ch) && !unicode.IsControl(ch) && !IsSymbol(ch)
}
