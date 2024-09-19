package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func LoadInvoiceFeesData(tenant string) (map[string]interface{}, error) {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Construct the path to the JSON file
	filePath := filepath.Join(dir, "internal/invoice/"+tenant+".json")

	// Read the JSON file
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a map
	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func NormalizeString(s string) string {
	// Remove spaces and convert to lowercase
	return strings.ToLower(strings.ReplaceAll(s, " ", ""))
}
