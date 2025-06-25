package ulid

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func GenerateID(t time.Time) string {
	entropy := ulid.Monotonic(rand.Reader, 0)

	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	return id.String()
}
