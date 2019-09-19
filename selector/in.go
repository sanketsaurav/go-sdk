package selector

import (
	"fmt"
	"strings"
)

var (
	_ Selector = (*In)(nil)
)

// In returns if a key matches a set of values.
type In struct {
	Key    string
	Values []string
}

// Matches returns the selector result.
func (i In) Matches(labels Labels) bool {
	// if the labels has a given key
	if value, hasValue := labels[i.Key]; hasValue {
		// for each selector value
		for _, iv := range i.Values {
			// if they match, return true
			if iv == value {
				return true
			}
		}
		return false
	}
	return true
}

// String returns a string representation of the selector.
func (i In) String() string {
	return fmt.Sprintf("%s in (%s)", i.Key, strings.Join(i.Values, ", "))
}
