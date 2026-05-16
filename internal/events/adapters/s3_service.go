package adapters

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Service struct {
	s3Client *s3.S3
	bucket   string
}

func NewS3Service(region, endpoint, bucket, accessKey, secretKey string, usePathStyle, useSSL bool) *s3Service {
	credentials := credentials.NewStaticCredentials(accessKey, secretKey, "")

	s3Session := session.Must(session.NewSession(&aws.Config{
		Region:           new(region),
		Endpoint:         new(endpoint),
		Credentials:      credentials,
		DisableSSL:       new(!useSSL),
		S3ForcePathStyle: new(usePathStyle),
	}))

	return &s3Service{
		s3Client: s3.New(s3Session),
		bucket:   bucket,
	}
}

func (s *s3Service) GetPresignedUploadURL(ctx context.Context, key, contentType string, expiresIn time.Duration) (string, error) {
	req, _ := s.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      new(s.bucket),
		Key:         new(key),
		ContentType: new(contentType),
	})

	url, err := req.Presign(expiresIn)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *s3Service) GetPresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: new(s.bucket),
		Key:    new(key),
	})

	url, err := req.Presign(expiresIn)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *s3Service) DeleteObject(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: new(s.bucket),
		Key:    new(key),
	}

	_, err := s.s3Client.DeleteObjectWithContext(ctx, input)
	return err
}
