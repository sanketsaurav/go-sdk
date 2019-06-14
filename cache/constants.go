package cache

// RemovalReason is a reason for removal.
type RemovalReason int

// RemovalReasons
const (
	Expired RemovalReason = iota
	Removed RemovalReason = iota
)
