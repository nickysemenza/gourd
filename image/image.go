package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/buckket/go-blurhash"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const aspectWidth = 4 * 2
const aspectHeight = 3 * 2

func GetBlurHash(ctx context.Context, url string) (hash string, image image.Image, err error) {
	ctx, span := otel.Tracer("image").Start(ctx, "image.GetBlurHash")
	defer span.End()

	image, err = GetFromURL(ctx, url)
	if err != nil {
		return
	}

	hash, err = blurhash.Encode(aspectWidth, aspectHeight, image)
	if err != nil {
		return
	}
	return
}

func GetFromURL(ctx context.Context, url string) (image.Image, error) {
	ctx, span := otel.Tracer("image").Start(ctx, "image.GetFromURL")
	defer span.End()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get image %s %w:", url, err)
	}
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get image %s %w:", url, err)
	}
	defer resp.Body.Close()

	var img io.Reader
	if strings.Contains(url, ".heic") {
		// shell out to imagemagick for heic conversion
		inputFile, err := ioutil.TempFile("", ".*.heic")
		if err != nil {
			return nil, err
		}
		defer inputFile.Close()
		outputFile, err := ioutil.TempFile("", ".*.png")
		if err != nil {
			return nil, err
		}
		defer outputFile.Close()

		_, err = io.Copy(inputFile, resp.Body)
		if err != nil {
			return nil, err
		}
		// nolint:gosec
		cmd := exec.Command("convert", inputFile.Name(), outputFile.Name())
		err = cmd.Start()
		if err != nil {
			return nil, err
		}
		err = cmd.Wait()
		if err != nil {
			return nil, err
		}
		rawImg, err := ioutil.ReadFile(outputFile.Name())
		if err != nil {
			return nil, err
		}
		img = bytes.NewReader(rawImg)
	} else {
		// jpg, png, gif
		img = resp.Body
	}

	m, format, err := image.Decode(img)
	if err != nil {
		return nil, fmt.Errorf("failed to get image %s %w:", url, err)
	}
	span.SetAttributes(attribute.String("image.format", format))
	return m, err

}
