package selector

import "fmt"

var (
	_ Selector = (*NotHasKey)(nil)
)

// NotHasKey returns if a label set does not have a given key.
type NotHasKey string

// Matches returns the selector result.
func (nhk NotHasKey) Matches(labels Labels) bool {
	if _, hasKey := labels[string(nhk)]; hasKey {
		return false
	}
	return true
}

// Validate validates the selector.
func (nhk NotHasKey) Validate(vr ValidationRules) (err error) {
	err = vr.CheckKey(string(nhk))
	return
}

// String returns a string representation of the selector.
func (nhk NotHasKey) String() string {
	return fmt.Sprintf("!%s", string(nhk))
}
