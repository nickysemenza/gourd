package handler

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	opentracing "github.com/opentracing/opentracing-go"
)

func getAWSSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		//LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
}

//AddFileToS3 puts a local file into s3 at a given path.
//Files are public with WRT their ACL.
func AddFileToS3(ctx context.Context, fileDir string, s3Path string) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "AddFileToS3")
	defer span.Finish()
	s, err := getAWSSession()
	if err != nil {
		log.Fatal(err)
	}

	span.LogEvent("opening file")
	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	span.LogEvent("reading file into memory")
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	log.Printf("saving file to S3 at %s", s3Path)
	put := s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("S3_BUCKET")),
		Key:           aws.String(s3Path),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(size),
		ACL:           aws.String("public-read"),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	}
	span.LogEventWithPayload("PutObject", put)

	resp, err := s3.New(s).PutObject(&put)
	span.LogEventWithPayload("finished", resp)
	return err
}

func ReadAndHash(r io.Reader) (io.Reader, string, error) {
	var b bytes.Buffer

	hash := md5.New()
	_, err := io.Copy(&b, io.TeeReader(r, hash))

	if err != nil {
		return nil, "", err
	}

	return &b, hex.EncodeToString(hash.Sum(nil)), nil
}

func GetImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
