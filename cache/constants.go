package cache

// RemovalReason is a reason for removal.
type RemovalReason int

// RemovalReasons
const (
	ExpiredTTL RemovalReason = iota
	Removed    RemovalReason = iota
)
