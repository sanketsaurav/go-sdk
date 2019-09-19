package selector

import "fmt"

var (
	_ Selector = (*NotEquals)(nil)
)

// NotEquals returns if a key strictly equals a value.
type NotEquals struct {
	Key, Value string
}

// Matches returns the selector result.
func (ne NotEquals) Matches(labels Labels) bool {
	if value, hasValue := labels[ne.Key]; hasValue {
		return ne.Value != value
	}
	return true
}

// String returns a string representation of the selector.
func (ne NotEquals) String() string {
	return fmt.Sprintf("%s != %s", ne.Key, ne.Value)
}
