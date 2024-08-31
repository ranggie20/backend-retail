package discard

import (
	"errors"
	"io"
)

type Bucket struct{}

func (b *Bucket) Upload(filename string, file io.Reader) (string, error) {
	return "", errors.New("discarding file: " + filename)
}
