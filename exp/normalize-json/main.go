package main

import (
	"encoding/json"
	"fmt"
	"sort"
)

func jsonNormalize(jsonStr string) ([]byte, error) {
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
		return nil, fmt.Errorf("JsonNormalize: json unmarshal: %w", err)
	}

	normalizedObj := normalize(jsonObj)

	normalizedJson, err := json.Marshal(normalizedObj)
	if err != nil {
		return nil, fmt.Errorf("JsonNormalize: json marshal: %w", err)
	}

	return normalizedJson, nil
}

func normalize(data interface{}) interface{} {
	switch v := data.(type) {
	case []interface{}:
		// Sort each element of the slice if it's a slice of maps or slices
		for i := range v {
			v[i] = normalize(v[i])
		}
		sort.Slice(v, func(i, j int) bool {
			return fmt.Sprintf("%v", v[i]) < fmt.Sprintf("%v", v[j])
		})
		return v
	case map[string]interface{}:
		// Recursively normalize each value in the map
		for key, value := range v {
			v[key] = normalize(value)
		}
	}

	return data
}

func main() {
	json1 := "{\"Id\":\"test-1\",\"Status\":\"Enabled\",\"Tierings\":[{\"AccessTier\":\"ARCHIVE_ACCESS\",\"Days\":120}],\"Filter\":{\"And\":{\"Prefix\":\"test\",\"Tags\":[{\"Key\":\"intelligent-Tiering\",\"Value\":\"true\"}]},\"Prefix\":null,\"Tag\":null}}\n"
	fmt.Printf("json1: %s\n", json1)

	normalizedJson1, err := jsonNormalize(string(json1))
	if err != nil {
		fmt.Println("Error normalizing JSON 1:", err)
		return
	}

	fmt.Printf("normalizedJson1: %s\n", normalizedJson1)
}
