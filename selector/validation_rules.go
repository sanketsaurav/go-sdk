package selector

// ValidationRules are rules used to validate selectors.
type ValidationRules interface {
	CheckKey(string) error
	CheckValue(string) error
	IsNameSymbol(rune) bool
}
