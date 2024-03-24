package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// Plug represents the structure of the plug block in the HCL config.
type Plug struct {
	Name       string     `hcl:"name,label"`
	Profile    string     `hcl:"profile"`
	Properties []Property `hcl:"property,block"`
}

// Property represents the nested property blocks within a plug block.
type Property struct {
	Type     string   `hcl:"type,label"`
	Name     string   `hcl:"name,label"`
	Compare  bool     `hcl:"compare"`
	Required bool     `hcl:"required"`
	Remain   hcl.Body `hcl:",remain"`
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
	diag = gohcl.DecodeBody(hclFile.Body, &hcl.EvalContext{}, &config)
	if diag.HasErrors() {
		fmt.Println("Error decoding HCL into Go struct:", diag)
		os.Exit(1)
	}

	for _, plug := range config.Plugs {
		fmt.Printf("Plug Name: %s\n", plug.Name)
		fmt.Printf("Profile: %s\n", plug.Profile)
		for _, property := range plug.Properties {
			fmt.Printf("Property Type: %s\n", property.Type)
			fmt.Printf("Property Name: %s\n", property.Name)
			fmt.Printf("Property Compare: %t\n", property.Compare)
			fmt.Printf("Property Required: %t\n", property.Required)

			if property.Remain != nil {
				// Attempt to decode the remaining body into a generic map
				var remain map[string]interface{}
				diag := gohcl.DecodeBody(property.Remain, &hcl.EvalContext{}, &remain)

				if diag.HasErrors() {
					fmt.Println("Error decoding remaining properties:", diag)
				} else if len(remain) == 0 {
					fmt.Println("No additional remain properties")
				} else {
					fmt.Printf("Additional remain properties: %+v\n", remain)
				}
			}
		}
	}

	// Output the decoded config to verify
	fmt.Printf("Parsed Config: %+v\n", config)
}
