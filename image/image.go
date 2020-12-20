package image

import (
	"context"
	"fmt"
	"image"
	"net/http"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/buckket/go-blurhash"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

const aspectWidth = 4
const aspectHeight = 3

func GetBlurHash(ctx context.Context, url string) (string, error) {
	ctx, span := otel.Tracer("image").Start(ctx, "image.GetBlurHash")
	defer span.End()

	image, err := GetFromURL(ctx, url)
	if err != nil {
		return "", err
	}

	return blurhash.Encode(aspectWidth, aspectHeight, image)
}

func GetFromURL(ctx context.Context, url string) (image.Image, error) {
	ctx, span := otel.Tracer("image").Start(ctx, "image.GetFromURL")
	defer span.End()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get image %s %w:", url, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get image %s %w:", url, err)
	}
	defer resp.Body.Close()

	m, format, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to get image %s %w:", url, err)
	}
	span.SetAttributes(label.String("image.format", format))
	return m, err

}
