package selector

import "fmt"

var (
	_ Selector = (*Equals)(nil)
)

// Equals returns if a key strictly equals a value.
type Equals struct {
	Key, Value string
}

// Matches returns the selector result.
func (e Equals) Matches(labels Labels) bool {
	if value, hasValue := labels[e.Key]; hasValue {
		return e.Value == value
	}
	return false
}

// String returns the string representation of the selector.
func (e Equals) String() string {
	return fmt.Sprintf("%s == %s", e.Key, e.Value)
}
