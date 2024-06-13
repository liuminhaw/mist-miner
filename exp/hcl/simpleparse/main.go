package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type HclConfig struct {
	Plugs []Plug `hcl:"plug,block"`
}

type Plug struct {
	Name          string            `hcl:"name,label"`
	Group         string            `hcl:"group,label"`
	Authenticator map[string]string `hcl:"authenticator,attr"`
	Diaries       []PlugDiary       `hcl:"diary,block"`
	Equipments    []PlugEquipment   `hcl:"equipment,block"`
}

type PlugEquipment struct {
	Type       string            `hcl:"type,label"`
	Name       string            `hcl:"name,label"`
	Attributes map[string]string `hcl:"attributes,attr"`
}

type PlugDiary struct {
	Type     string `hcl:"type,label"`
	Name     string `hcl:"name,label"`
	Compare  bool   `hcl:"compare,optional"`
	Required bool   `hcl:"required,optional"`
}

func main() {
	// Initialize HCL parser
	parser := hclparse.NewParser()

	// Parse the HCL file
	hclFile, diag := parser.ParseHCLFile("config.hcl")
	if diag.HasErrors() {
		fmt.Println("Error parsing HCL file:", diag)
		os.Exit(1)
	}

	// Decode the file into Go struct
	var config struct {
		Plugs []Plug `hcl:"plug,block"`
	}
	diag = gohcl.DecodeBody(hclFile.Body, nil, &config)
	if diag.HasErrors() {
		fmt.Println("Error decoding HCL into Go struct:", diag)
		os.Exit(1)
	}

	for _, plug := range config.Plugs {
		fmt.Printf("Plug Name: %s\n", plug.Name)
		fmt.Printf("Plug Group: %s\n", plug.Group)

		// Print the authenticator
		for key, value := range plug.Authenticator {
			fmt.Printf("%s: %s\n", key, value)
		}

		// Print the diaries
		for _, diary := range plug.Diaries {
			fmt.Printf("Diary Type: %s\n", diary.Type)
			fmt.Printf("Diary Name: %s\n", diary.Name)
			fmt.Printf("Diary Compare: %t\n", diary.Compare)
			fmt.Printf("Diary Required: %t\n", diary.Required)
		}

        // Print the equipments
        for _, equipment := range plug.Equipments {
            fmt.Printf("Equipment Type: %s\n", equipment.Type)
            fmt.Printf("Equipment Name: %s\n", equipment.Name)

            for key, value := range equipment.Attributes {
                fmt.Printf("%s: %s\n", key, value)
            }
        }

		// for _, property := range plug.Properties {
		// 	fmt.Printf("Property Type: %s\n", property.Type)
		// 	fmt.Printf("Property Name: %s\n", property.Name)
		// 	fmt.Printf("Property Compare: %t\n", property.Compare)
		// 	fmt.Printf("Property Required: %t\n", property.Required)
		//
		// 	if property.Remain != nil {
		// 		// Attempt to decode the remaining body into a generic map
		// 		var remain map[string]interface{}
		// 		diag := gohcl.DecodeBody(property.Remain, &hcl.EvalContext{}, &remain)
		//
		// 		if diag.HasErrors() {
		// 			fmt.Println("Error decoding remaining properties:", diag)
		// 		} else if len(remain) == 0 {
		// 			fmt.Println("No additional remain properties")
		// 		} else {
		// 			fmt.Printf("Additional remain properties: %+v\n", remain)
		// 		}
		// 	}
		// }
	}

	// Output the decoded config to verify
	// fmt.Printf("Parsed Config: %+v\n", config)
}
