package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// RecipeFromFile reads a recipe from json file
func RecipeFromFile(ctx context.Context, inputPath string) (*RecipeWrapper, error) {
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

	switch filepath.Ext(inputPath) {
	case ".json":
	case ".yaml":
		y, err := yaml.YAMLToJSON(fileBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to read recipe: %w", err)
		}
		fileBytes = y
	default:
		return nil, fmt.Errorf("unknown extension: %s", inputPath)
	}

	// we initialize our Users array
	var r RecipeWrapper

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(fileBytes, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe: %w", err)
	}

	return &r, nil
}
