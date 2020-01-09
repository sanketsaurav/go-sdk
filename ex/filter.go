package ex

// Filter returns a predicate set for a given
// set of predicates.
func Filter(preds ...Predicate) Predicates {
	return Predicates(preds)
}

// FilterIs returns a Predicate that wraps
// ex.Is(...)
func FilterIs(cause error) Predicate {
	return func(err error) bool {
		return Is(err, cause)
	}
}

// Predicate is a predicate for an exception.
type Predicate func(error) bool

// Predicates are a set of filters, if any mat
type Predicates []Predicate

// Any applies the filters to a given error,
// returning true if any filters pass.
func (p Predicates) Any(err error) bool {
	for _, filter := range p {
		if filter(err) {
			return true
		}
	}
	return false
}

// All applies the filters to a given error,
// returning true if all filters pass.
func (p Predicates) All(err error) bool {
	for _, filter := range p {
		if !filter(err) {
			return false
		}
	}
	return true
}
