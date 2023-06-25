package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// func processLine[T any](line string) {
// 	var data T
// 	err := json.Unmarshal([]byte(line), &data)
// 	if err != nil {
// 		fmt.Printf("Error unmarshaling JSON: %v\n", err)
// 		return
// 	}

// 	// Process the data from the line as desired
// 	// fmt.Printf("ok")
// 	// ...
// }

// func loadNDJSONFile[T any](filename string) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		processLine[T](line)
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return err
// 	}

// 	return nil
// }

func processLines[T any](lines []string) {
	for _, line := range lines {
		var data T
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			continue
		}

		// Process the data from the line as desired
		// ...
	}
	fmt.Printf("did batch of %d\n", len(lines))
}

func loadNDJSONFile[T any](filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0, 1000)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		if len(lines) == 1000 {
			processLines[T](lines)
			lines = make([]string, 0, 1000)
		}
	}

	if len(lines) > 0 {
		processLines[T](lines)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func Test() {
	filename := "/Users/nicky/dev/gourd/tmp/usda_json/brandedfoods.ndjson"

	err := loadNDJSONFile[BrandedFoodItem](filename)
	if err != nil {
		log.Fatal(err)
	}
}
