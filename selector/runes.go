package selector

// IsSymbol returns if a rune is a symbol.
func IsSymbol(ch rune) bool {
	return (int(ch) >= int(Bang) && int(ch) <= int(ForwardSlash)) ||
		(int(ch) >= int(Colon) && int(ch) <= int(At)) ||
		(int(ch) >= int(OpenBracket) && int(ch) <= int(BackTick)) ||
		(int(ch) >= int(OpenCurly) && int(ch) <= int(Tilde))
}

// IsNameSymbol returns if a rune is a name symbol.
func IsNameSymbol(ch rune) bool {
	switch ch {
	case Dot, Dash, Underscore:
		return true
	}
	return false
}

// IsSelectorSymbol returns if a rune is a selector symbol.
func IsSelectorSymbol(ch rune) bool {
	switch ch {
	case Equal, Bang, OpenParens, CloseParens, Comma:
		return true
	}
	return false
}
