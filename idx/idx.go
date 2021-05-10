package idx

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"regexp"

	"github.com/google/uuid"
)

var (
	pattern = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
)

func NewUUID() string {
	return uuid.New().String()
}

// IsUUIDString は文字列が UUID かどうか判定する。
func IsUUIDString(s string) bool {
	return pattern.MatchString(s)
}

// ShortID produces a "unique" 6 bytes long string.
// Do not use as a reliable way to get unique IDs, instead use for things like logging.
func ShortID() string {
	b := make([]byte, 6)
	_, _ = io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}
