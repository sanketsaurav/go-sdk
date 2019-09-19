package selector

// RuneSet takes a list of runes and returns a lookup.
func RuneSet(runes ...rune) map[rune]struct{} {
	output := make(map[rune]struct{})
	for _, r := range runes {
		output[r] = struct{}{}
	}
	return output
}
