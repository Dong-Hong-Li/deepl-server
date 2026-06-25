package utils

import (
	"crypto/rand"

	"github.com/oklog/ulid/v2"
)

func NewFileKey(prefix string, ext string) string {
	return prefix + "_" + ulid.MustNew(
		ulid.Now(),
		rand.Reader,
	).String() + "." + ext
}
