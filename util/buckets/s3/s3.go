package s3

import (
	"fmt"
	"io"

	"github.com/online-bnsp/backend/util/buckets"
	"github.com/online-bnsp/backend/util/s3"
)

type Bucket struct {
	s3.Bucket
}

func New(bucket s3.Bucket) buckets.Bucket {
	return &Bucket{Bucket: bucket}
}

func (b *Bucket) Upload(filename string, file io.Reader) (string, error) {
	err := b.Bucket.Connect()
	if err != nil {
		return "", fmt.Errorf("connect error: %w", err)
	}
	return b.Bucket.Upload(filename, file)
}
