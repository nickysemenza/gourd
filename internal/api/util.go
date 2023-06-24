package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ericlagergren/decimal"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/clients/rs_client"
	"github.com/nickysemenza/gourd/internal/common"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/yaml"
)

func l(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx)
}

func parsePagination(o *OffsetParam, l *LimitParam) Items {
	offset := 0
	limit := 20
	if o != nil {
		offset = int(*o)
	}
	if l != nil {
		limit = int(*l)
	}
	items := Items{
		Offset: offset,
		Limit:  limit,
	}
	if limit == 0 {
		items.PageNumber = 0
	} else {
		items.PageNumber = offset/limit + 1
	}

	return items
}

func (l *Items) setTotalCount(count uint64) {
	c := int(count)
	l.TotalCount = c
	l.PageCount = int(math.Ceil(float64(c) / float64(l.Limit)))
}

func sendErr(c echo.Context, code int, err error) error {
	ctx := c.Request().Context()
	trace.SpanFromContext(ctx).AddEvent(fmt.Sprintf("error: %v", err))
	logrus.WithField("code", code).WithField("route", c.Request().URL).Errorf("http err: %v", err)
	t := GetTraceID(ctx)
	return c.JSON(code, Error{
		Message: err.Error(),
		TraceId: &t,
	})
}

func handleErr(c echo.Context, err error) error {
	if errors.Is(err, common.ErrNotFound) {
		return sendErr(c, http.StatusNotFound, err)
	}
	return sendErr(c, http.StatusInternalServerError, err)

}

func GetTraceID(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		return sc.TraceID().String()
	}
	return ""
}

func bytesFromFile(_ context.Context, inputPath string) ([]byte, error) {
	inputFile, err := os.Open(inputPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	fmt.Println("Successfully Opened", inputPath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer inputFile.Close()

	return io.ReadAll(inputFile)

}
func JSONBytesFromFile(ctx context.Context, inputPath string) ([][]byte, error) {
	fileBytes, err := bytesFromFile(ctx, inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	fileDocs := bytes.Split(fileBytes, []byte("---\n"))
	var output [][]byte
	for _, doc := range fileDocs {
		if len(doc) < len("recipe") {
			// too short to be anything
			continue
		}
		switch filepath.Ext(inputPath) {
		case ".json":
		case ".yaml":
			y, err := yaml.YAMLToJSON(doc)
			if err != nil {
				return nil, fmt.Errorf("failed to read yaml: %w", err)
			}
			doc = y
		default:
			return nil, fmt.Errorf("unknown extension: %s", inputPath)
		}

		output = append(output, doc)
	}

	return output, nil
}
func (a *API) RecipeFromText(ctx context.Context, text string) (*RecipeDetailInput, error) {
	ctx, span := a.tracer.Start(ctx, "RecipeFromText")
	defer span.End()

	output := RecipeDetailInput{}
	err := a.R.Call(ctx, text, rs_client.RecipeDecode, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to decode recipe: %w", err)
	}
	return &output, nil
}
func (a *API) RecipeFromCompact(ctx context.Context, cr CompactRecipe) (*RecipeWrapperInput, error) {
	output := RecipeWrapperInput{}
	err := a.R.Send(ctx, "codec/expand", cr, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to decode recipe: %w", err)
	}
	return &output, nil
}
func (a *API) NormalizeAmount(ctx context.Context, amt Amount) (*Amount, error) {
	output := Amount{}
	err := a.R.Send(ctx, "normalize_amount", amt, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize recipe: %w", err)
	}
	return &output, nil
}

// RecipeFromFile reads a recipe from json or yaml file
func (a *API) RecipeFromFile(ctx context.Context, inputPath string) (output []RecipeDetailInput, error error) {
	if strings.HasSuffix(inputPath, "/") {
		// import as directory
		dirEntries, err := os.ReadDir(inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read: %w", err)
		}
		for _, entry := range dirEntries {
			file := filepath.Join(inputPath, entry.Name())
			switch filepath.Ext(file) {
			case ".txt", ".json", ".yaml":
				recipe, err := a.RecipeFromFile(ctx, file)
				if err != nil {
					return nil, fmt.Errorf("failed to read: %w", err)
				}
				output = append(output, recipe...)
			default:
				continue
			}

		}
		return
	}
	l(ctx).Infof("loading %s", inputPath)

	switch filepath.Ext(inputPath) {
	case ".txt":
		data, err := bytesFromFile(ctx, inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read bytes: %w", err)
		}
		output, err := a.RecipeFromText(ctx, string(data))
		if err != nil {
			return nil, fmt.Errorf("failed to decode recipe: %w", err)
		}
		output.Sources = &[]RecipeSource{{Title: &inputPath}}
		return []RecipeDetailInput{*output}, nil
	case ".json", ".yaml":

		jsonBytes, err := JSONBytesFromFile(ctx, inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read recipe: %w", err)
		}

		for _, doc := range jsonBytes {

			var r RecipeDetailInput
			err = json.Unmarshal(doc, &r)
			if err != nil {
				return nil, fmt.Errorf("failed to read recipe: %w", err)
			}
			if r.Name == "" {
				return nil, fmt.Errorf("failed to read recipe name from file %s [%v]", inputPath, r)
			}
			r.Sources = &[]RecipeSource{{Title: &inputPath}}
			output = append(output, r)
		}
		return
	default:
		return nil, fmt.Errorf("unknown extension: %s", inputPath)
	}

}

// IngredientMappingFromFile is todo
func IngredientMappingFromFile(ctx context.Context, inputPath string) ([]IngredientMapping, error) {
	jsonBytes, err := JSONBytesFromFile(ctx, inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe: %w", err)
	}

	var r []IngredientMapping
	err = json.Unmarshal(jsonBytes[0], &r)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe: %w", err)
	}

	return r, nil
}

// FetchAndTransform returns a recipe.
func (a *API) FetchAndTransform(ctx context.Context, addr string, ingredientToId func(ctx context.Context, name string) (string, error)) (*RecipeWrapperInput, error) {
	ctx, span := otel.Tracer("scraper").Start(ctx, "scraper.GetIngredients")
	defer span.End()

	r := RecipeWrapperInput{}
	err := a.R.Call(ctx, addr, rs_client.Scrape, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Coalesce returns the first non-zero provided
func Coalesce[T string | int | float32 | float64](i ...T) T {
	var zeroVal T
	for _, s := range i {
		if s != zeroVal {
			return s
		}
	}
	return zeroVal
}

func decimalFromFloat(f float64) types.Decimal {
	d := decimal.WithContext(types.DecimalContext)
	d.SetFloat64(f)
	return types.NewDecimal(d)
}
func nullDecimalFromFloat(f *float64) types.NullDecimal {
	if f == nil {
		return types.NewNullDecimal(nil)
	}
	d := decimal.WithContext(types.DecimalContext)
	d.SetFloat64(*f)
	return types.NewNullDecimal(d)
}
