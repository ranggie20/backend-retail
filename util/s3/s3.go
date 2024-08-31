package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Bucket interface {
	Connect() error
	Upload(filename string, file io.Reader) (string, error)
}

type s3 struct {
	uri        string
	accessKey  string
	secretKey  string
	token      string
	region     string
	bucketName string
	session    *session.Session
}

func NewS3(uri, accessKey, secretKey, token, region, bucketName string) Bucket {
	return &s3{
		uri:        uri,
		accessKey:  accessKey,
		secretKey:  secretKey,
		token:      token,
		region:     region,
		bucketName: bucketName,
	}
}

func (b *s3) Connect() error {
	sess, err := session.NewSession(
		&aws.Config{
			Endpoint: &b.uri,
			Region:   aws.String(b.region),
			Credentials: credentials.NewStaticCredentials(
				b.accessKey,
				b.secretKey,
				b.token, // a token will be created when the session it's used.
			),
		})
	if err != nil {
		return err
	}
	b.session = sess

	return nil
}

func (b *s3) Upload(filename string, file io.Reader) (string, error) {
	uploader := s3manager.NewUploader(b.session)

	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(b.bucketName),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	return up.Location, nil
}
