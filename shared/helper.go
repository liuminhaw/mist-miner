package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// JsonNormalize normalizes a JSON string by sorting the keys of each object in the JSON string.
func JsonNormalize(jsonStr string) ([]byte, error) {
	// fmt.Printf("JsonNormalize: %s\n", jsonStr)

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

// JsonMarshal marshals a JSON object to a byte slice.
func JsonMarshal(data any) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("JsonMarshal: %w", err)
	}

	return buffer.Bytes(), nil
}

// normalize recursively normalizes a JSON object by sorting the keys of each object.
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
