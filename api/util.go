package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// RecipeFromFile reads a recipe from json file
func RecipeFromFile(ctx context.Context, path string) (*RecipeDetail, error) {
	jsonFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("failed to read recipe: %w", err)
	}
	fmt.Println("Successfully Opened", path)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe: %w", err)
	}
	// we initialize our Users array
	var r RecipeDetail

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe: %w", err)
	}

	return &r, nil
}
