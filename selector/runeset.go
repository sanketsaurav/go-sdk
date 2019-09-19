package selector

// NewRuneSet takes a list of runes and returns a lookup.
func NewRuneSet(runes ...rune) RuneSet {
	output := make(map[rune]struct{})
	for _, r := range runes {
		output[r] = struct{}{}
	}
	return output
}

// RuneSet is a lookup for runes.
type RuneSet map[rune]struct{}

// Has returns if a given rune is in the runset.
func (rs RuneSet) Has(ch rune) bool {
	_, ok := rs[ch]
	return ok
}
