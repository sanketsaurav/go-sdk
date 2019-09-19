package selector

var (
	_ Selector = (*HasKey)(nil)
)

// HasKey returns if a label set has a given key.
type HasKey string

// Matches returns the selector result.
func (hk HasKey) Matches(labels Labels) bool {
	_, hasKey := labels[string(hk)]
	return hasKey
}

// String returns a string representation of the selector.
func (hk HasKey) String() string {
	return string(hk)
}
