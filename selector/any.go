package selector

var (
	_ Selector = (*Any)(nil)
)

// Any matches everything
type Any struct{}

// Matches returns true
func (a Any) Matches(labels Labels) bool {
	return true
}

// String returns a string representation of the selector
func (a Any) String() string {
	return ""
}
