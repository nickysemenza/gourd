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
	"sigs.k8s.io/yaml"
)

func bytesFromFile(ctx context.Context, inputPath string) ([]byte, error) {
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

// RecipeFromFile reads a recipe from json or yaml file
func RecipeFromFile(ctx context.Context, inputPath string) ([]RecipeDetail, error) {
	ext := filepath.Ext(inputPath)
	log.Infof("loading %s (%s)", inputPath, ext)
	switch ext {
	case ".txt":
		data, err := bytesFromFile(ctx, inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read bytes: %w", err)
		}
		output := RecipeDetail{}
		err = rs_client.Call(ctx, string(data), rs_client.RecipeDecode, &output)
		if err != nil {
			return nil, fmt.Errorf("failed to decode recipe: %w", err)
		}
		return []RecipeDetail{output}, nil
	case ".json", ".yaml":
		var output []RecipeDetail

		jsonBytes, err := JSONBytesFromFile(ctx, inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read recipe: %w", err)
		}

		for _, doc := range jsonBytes {

			var r RecipeDetail
			err = json.Unmarshal(doc, &r)
			if err != nil {
				return nil, fmt.Errorf("failed to read recipe: %w", err)
			}
			if r.Name == "" {
				return nil, fmt.Errorf("failed to read recipe name from file %s [%v]", inputPath, r)
			}
			output = append(output, r)
		}
		return output, nil
	default:
		return nil, fmt.Errorf("unknown extension: %s", inputPath)
	}
	return nil, nil
}

type IngredientMapping struct {
	Name         string        `json:"name"`
	FdcID        int           `json:"fdc_id"`
	Aliases      []string      `json:"aliases"`
	UnitMappings []UnitMapping `json:"unit_mappings"`
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
