package buckets

import "io"

type Bucket interface {
	Upload(filename string, file io.Reader) (string, error)
}
