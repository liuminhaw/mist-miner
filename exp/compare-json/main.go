package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

// jsonEqual compares two JSON strings considering complex nested structures
func jsonEqual(json1, json2 string) (bool, error) {
	var obj1 interface{}
	var obj2 interface{}

	err := json.Unmarshal([]byte(json1), &obj1)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(json2), &obj2)
	if err != nil {
		return false, err
	}

	// Normalize the objects (sorting arrays, etc.)
	normalizedObj1 := normalize(obj1)
	normalizedObj2 := normalize(obj2)

    fmt.Printf("normalizedObj1: %+v\n", normalizedObj1)
    fmt.Printf("normalizedObj2: %+v\n", normalizedObj2)

    normalizedJson1, err := json.Marshal(normalizedObj1)
    if err != nil {
        return false, err
    }
    normalizedJson2, err := json.Marshal(normalizedObj2)
    if err != nil {
        return false, err
    }
    fmt.Printf("normalizedJson1: %s\n", normalizedJson1)
    fmt.Printf("normalizedJson2: %s\n", normalizedJson2)

	// Using reflect.DeepEqual to compare the normalized objects
	return reflect.DeepEqual(normalizedObj1, normalizedObj2), nil
}

// normalize recursively sorts any arrays it encounters
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
	json1 := `{"Id":"test-1","Status":"Enabled","Tierings":[{"AccessTier":"ARCHIVE_ACCESS","Days":120}],"Filter":{"And":{"Prefix":"pre-test","Tags":[{"Key":"intelligent-Tiering","Value":"true"}]},"Prefix":null,"Tag":null}}`
	json2 := `{"Status":"Enabled","Id":"test-1","Tierings":[{"AccessTier":"ARCHIVE_ACCESS","Days":120}],"Filter":{"And":{"Prefix":"pre-test","Tags":[{"Key":"intelligent-Tiering","Value":"true"}]},"Prefix":null,"Tag":null}}`

	equal, err := jsonEqual(json1, json2)
	if err != nil {
		fmt.Println("Error comparing JSON:", err)
		return
	}

	fmt.Println("Are JSON strings equal?:", equal)
}
