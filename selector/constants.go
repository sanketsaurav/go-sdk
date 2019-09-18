package selector

const (
	// At is a common rune.
	At = rune('@')
	// Colon is a common rune.
	Colon = rune(':')
	// Dash is a common rune.
	Dash = rune('-')
	// Underscore  is a common rune.
	Underscore = rune('_')
	// Dot is a common rune.
	Dot = rune('.')
	// ForwardSlash is a common rune.
	ForwardSlash = rune('/')
	// BackSlash is a common rune.
	BackSlash = rune('\\')
	// BackTick is a common rune.
	BackTick = rune('`')
	// Bang is a common rune.
	Bang = rune('!')
	// Comma is a common rune.
	Comma = rune(',')
	// OpenBracket is a common rune.
	OpenBracket = rune('[')
	// OpenParens is a common rune.
	OpenParens = rune('(')
	// OpenCurly is a common rune.
	OpenCurly = rune('{')
	// CloseBracket is a common rune.
	CloseBracket = rune(']')
	// CloseParens is a common rune.
	CloseParens = rune(')')
	// Equal is a common rune.
	Equal = rune('=')
	// Space is a common rune.
	Space = rune(' ')
	// Tab is a common rune.
	Tab = rune('\t')
	// Tilde is a common rune.
	Tilde = rune('~')
	// CarriageReturn is a common rune.
	CarriageReturn = rune('\r')
	// NewLine is a common rune.
	NewLine = rune('\n')
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
	ErrKeyDNSPrefixTooLong = Error("key dns prefix too long; must be less than 253 characters")
	// ErrValueTooLong indicates a value is too long.
	ErrValueTooLong = Error("value too long; must be less than 63 characters")
	// ErrKeyInvalidCharacter indicates a key contains characters
	ErrKeyInvalidCharacter = Error(`key contains invalid characters, regex used: ([A-Za-z0-9_-\.])`)
)

var (
	// DefaultValidationRules are the default validation options.
	DefaultValidationRules = Kubernetes{
		// MaxDNSPrefixLen is the maximum dns prefix length.
		MaxDNSPrefixLen: 253,
		// MaxKeyLen is the maximum key length.
		MaxKeyLen: 63,
		// MaxValueLen is the maximum value length.
		MaxValueLen: 63,
	}
)
