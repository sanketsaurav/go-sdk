package selector

// Option is a tweak to selector parsing.
type Option func(po *Parser)

// OptValidationRules is an option to set the validation rules for the parser.
func OptValidationRules(vr ValidationRules) Option {
	return func(p *Parser) {
		p.ValidationRules = vr
	}
}

// OptSkipValidation is an option to skip checking the values of selector expressions.
func OptSkipValidation() Option {
	return func(p *Parser) {
		p.SkipValidation = true
	}
}
