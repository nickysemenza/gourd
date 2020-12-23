package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// RecipeFromFile reads a recipe from json or yaml file
func RecipeFromFile(ctx context.Context, inputPath string) ([]RecipeDetail, error) {
	inputFile, err := os.Open(inputPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read recipe: %w", err)
	}
	fmt.Println("Successfully Opened", inputPath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer inputFile.Close()

	fileBytes, err := ioutil.ReadAll(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe: %w", err)
	}
	fileDocs := bytes.Split(fileBytes, []byte("---\n"))
	var output []RecipeDetail
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
				return nil, fmt.Errorf("failed to read recipe: %w", err)
			}
			doc = y
		default:
			return nil, fmt.Errorf("unknown extension: %s", inputPath)
		}

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
}
