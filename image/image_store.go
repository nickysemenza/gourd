package image

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Store interface {
	GetImageURL(ctx context.Context, id string) string
	SaveImage(ctx context.Context, id string, data image.Image) error
	Dir() string
}

type LocalImageStore struct {
	httpBaseURL, dir string
	tracer           trace.Tracer
}

func (l *LocalImageStore) Dir() string {
	return l.dir
}

func NewLocalImageStore(httpBaseURL string) (*LocalImageStore, error) {
	dir := "/tmp/gourd_images"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return &LocalImageStore{
		httpBaseURL: httpBaseURL,
		dir:         dir,
		tracer:      otel.Tracer("db"),
	}, nil
}

func (l *LocalImageStore) GetImageURL(ctx context.Context, id string) string {
	fileName := l.getFileName(id)
	return l.httpBaseURL + "/images/" + fileName
}
func (l *LocalImageStore) getFileName(id string) string {
	return fmt.Sprintf("%s.png", id)
}
func (l *LocalImageStore) SaveImage(ctx context.Context, id string, data image.Image) error {
	_, span := l.tracer.Start(ctx, "SaveImage")
	defer span.End()
	fileName := l.getFileName(id)
	fileName = l.dir + "/" + fileName
	outputFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, data)
	if err != nil {
		return err
	}
	logrus.Infof("Saved image %s", fileName)

	return nil
}
