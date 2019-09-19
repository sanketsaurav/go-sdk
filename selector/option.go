package selector

// Option is a tweak to selector parsing.
type Option func(*Parser)

// OptSkipValidation disables all validation.
func OptSkipValidation() Option {
	return func(p *Parser) {
		p.SkipValidation = true
	}
}

// OptAllowSymbolRepeats allows symbols to repeat.
func OptAllowSymbolRepeats() Option {
	return func(p *Parser) {
		p.AllowSymbolRepeats = true
	}
}

// OptAllowSymbolPrefix allows symbol prefixes for names.
func OptAllowSymbolPrefix() Option {
	return func(p *Parser) {
		p.AllowSymbolPrefix = true
	}
}

// OptAllowSymbolSuffix allows symbol suffixes for names.
func OptAllowSymbolSuffix() Option {
	return func(p *Parser) {
		p.AllowSymbolSuffix = true
	}
}

// OptNameSymbols sets the allowed name symbols.
func OptNameSymbols(runes ...rune) Option {
	return func(p *Parser) {
		p.NameSymbols = RuneSet(runes...)
	}
}
