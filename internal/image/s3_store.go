//nolint:staticcheck
package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type S3Store struct {
	tracer         trace.Tracer
	s3Client       *s3.Client
	presignClient  *s3.PresignClient
	bucket, prefix string
}

var _ Store = &S3Store{}

func NewS3Store(ctx context.Context, endpoint, region, bucket, keyID, appKey, prefix string) (*S3Store, error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(keyID, appKey, ""))
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint}, nil
			})),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })
	presignClient := s3.NewPresignClient(s3Client)

	return &S3Store{
		otel.Tracer("s3"),
		s3Client,
		presignClient,
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
	_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
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
	req, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.prefix + "/" + getFileName(id)),
	})

	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s *S3Store) Dir() string {
	return ""
}
