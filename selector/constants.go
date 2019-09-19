package selector

// Operators
const (
	OpEquals       = "="
	OpDoubleEquals = "=="
	OpNotEquals    = "!="
	OpIn           = "in"
	OpNotIn        = "notin"
)

// Common runes
const (
	At             = rune('@')
	BackSlash      = rune('\\')
	BackTick       = rune('`')
	Bang           = rune('!')
	CarriageReturn = rune('\r')
	CloseBracket   = rune(']')
	CloseParens    = rune(')')
	Colon          = rune(':')
	Comma          = rune(',')
	Dash           = rune('-')
	Dot            = rune('.')
	Equal          = rune('=')
	ForwardSlash   = rune('/')
	NewLine        = rune('\n')
	OpenBracket    = rune('[')
	OpenCurly      = rune('{')
	OpenParens     = rune('(')
	Space          = rune(' ')
	Star           = rune('*')
	Tab            = rune('\t')
	Tilde          = rune('~')
	Underscore     = rune('_')
)

// Validation defaults
const (
	DefaultMaxDNSPrefixLen = 253
	DefaultMaxKeyLen       = 63
	DefaultMaxValueLen     = 63
)

// Validation defaults
var (
	DefaultNameSymbols     = []rune{Dot, Dash, Underscore}
	DefaultSelectorSymbols = []rune{Equal, Bang, OpenParens, CloseParens, Comma}
)

const (
	// ErrInvalidOperator is returned if the operator is invalid.
	ErrInvalidOperator = Error("invalid operator")
	// ErrInvalidSelector is returned if there is a structural issue with the selector.
	ErrInvalidSelector = Error("invalid selector")
	// ErrKeyEmpty indicates a key is empty.
	ErrKeyEmpty = Error("key empty")
	// ErrKeyTooLong indicates a key is too long.
	ErrKeyTooLong = Error("key too long")
	// ErrKeyDNSPrefixEmpty indicates a key's "dns" prefix is empty.
	ErrKeyDNSPrefixEmpty = Error("key dns prefix empty")
	// ErrKeyDNSPrefixTooLong indicates a key's "dns" prefix is empty.
	ErrKeyDNSPrefixTooLong = Error("key dns prefix too long")
	// ErrValueTooLong indicates a value is too long.
	ErrValueTooLong = Error("value too long")
	// ErrKeyInvalidCharacter indicates a key contains characters
	ErrKeyInvalidCharacter = Error(`key or value contains invalid characters`) // , regex used: ([A-Za-z0-9_-\.])
)
