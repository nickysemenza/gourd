package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type S3Store struct {
	tracer         trace.Tracer
	s3Client       *s3.S3
	bucket, prefix string
}

var _ Store = &S3Store{}

func NewS3Store(endpoint, region, bucket, keyID, appKey, prefix string) (*S3Store, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(keyID, appKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}
	s3Client := s3.New(newSession)

	return &S3Store{
		otel.Tracer("s3"),
		s3Client,
		bucket,
		prefix,
	}, nil
}
func (s *S3Store) SaveImage(ctx context.Context, id string, data image.Image) error {
	ctx, span := s.tracer.Start(ctx, "SaveImage")
	defer span.End()
	buf := new(bytes.Buffer)
	err := png.Encode(buf, data)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("s3store: %w", err)
	}

	// Upload a new object "testfile.txt" with the string "S3 Compatible API"
	_, err = s.s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.prefix + "/" + getFileName(id)),
	})
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("s3store: %w", err)
	}
	return nil

}
func (s *S3Store) GetImageURL(ctx context.Context, id string) (string, error) {
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.prefix + "/" + getFileName(id)),
	})

	url, err := req.Presign(90 * time.Minute) // Set link expiration time
	if err != nil {
		return "", err
	}

	return url, err
}

func (s *S3Store) Dir() string {
	return ""
}
