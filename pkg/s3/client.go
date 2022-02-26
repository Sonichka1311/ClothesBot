package s3

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Endpoint = "https://storage.yandexcloud.net"
	bucketName = ""
)

type S3 struct {
	*s3.S3
	bucketName string
}

func NewS3() *S3 {
	s3Session, _ := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: &s3Endpoint,
	})

	return &S3{
		S3:         s3.New(s3Session),
		bucketName: bucketName,
	}
}

func (s S3) DeleteObject(key string) error {
	_, err := s.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	return err
}

func (s S3) GetObject(key string) (io.ReadCloser, error) {
	res, err := s.S3.GetObject(&s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func (s S3) PutObject(key string, data []byte) error {
	body := bytes.NewReader(data)
	_, err := s.S3.PutObject(&s3.PutObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
		Body:   body,
	})
	return err
}