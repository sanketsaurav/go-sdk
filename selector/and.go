package selector

import "strings"

var (
	_ Selector = (*And)(nil)
)

// And is a combination selector.
type And []Selector

// Matches returns if both A and B match the labels.
func (a And) Matches(labels Labels) bool {
	for _, s := range a {
		if !s.Matches(labels) {
			return false
		}
	}
	return true
}

// And returns a string representation for the selector.
func (a And) String() string {
	var childValues []string
	for _, c := range a {
		childValues = append(childValues, c.String())
	}
	return strings.Join(childValues, ", ")
}
