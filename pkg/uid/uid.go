package uid

import (
	"crypto/rand"
	"encoding/hex"
)

// New generates a random 24-char hex string compatible with MongoDB ObjectID format.
func New() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
