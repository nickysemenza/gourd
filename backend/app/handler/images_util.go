package handler

import (
	"bytes"
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
)

func getAWSSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		//LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
}

//AddFileToS3 puts a local file into s3 at a given path.
//Files are public with WRT their ACL.
func AddFileToS3(fileDir string, s3Path string) error {

	s, err := getAWSSession()
	if err != nil {
		log.Fatal(err)
	}

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	log.Printf("saving file to S3 at %s", s3Path)
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("S3_BUCKET")),
		Key:           aws.String(s3Path),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(size),
		ACL:           aws.String("public-read"),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	})
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
