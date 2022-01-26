package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nickysemenza/gourd/rs_client"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"sigs.k8s.io/yaml"
)

func bytesFromFile(_ context.Context, inputPath string) ([]byte, error) {
	inputFile, err := os.Open(inputPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	fmt.Println("Successfully Opened", inputPath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer inputFile.Close()

	return ioutil.ReadAll(inputFile)

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

// RecipeFromFile reads a recipe from json or yaml file
func (a *API) RecipeFromFile(ctx context.Context, inputPath string) ([]RecipeDetailInput, error) {
	ext := filepath.Ext(inputPath)
	log.Infof("loading %s (%s)", inputPath, ext)
	switch ext {
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
		var output []RecipeDetailInput

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
		return output, nil
	default:
		return nil, fmt.Errorf("unknown extension: %s", inputPath)
	}

}

// IngredientMapping is todo
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
func (a *API) FetchAndTransform(ctx context.Context, addr string, ingredientToId func(ctx context.Context, name string, kind string) (string, error)) (*RecipeWrapperInput, error) {
	ctx, span := otel.Tracer("scraper").Start(ctx, "scraper.GetIngredients")
	defer span.End()

	r := RecipeWrapperInput{}
	err := a.R.Call(ctx, addr, rs_client.Scrape, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
