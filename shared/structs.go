package shared

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type MinerConfig struct {
	Path string
}

type MinerData struct {
	Type  string
	Name  string
	Value string
}

type MinerResource []MinerData

type MinerResources []MinerResource

// HCL config structure
type HclConfig struct {
	Plugs []Plug `hcl:"plug,block"`
}

type Plug struct {
	Name       string     `hcl:"name,label"`
	Profile    string     `hcl:"profile"`
	Identifier Identifier `hcl:"identifier,block"`
	Properties []Property `hcl:"property,block"`
}

type Identifier struct {
	Field string `hcl:"field"`
}

type Property struct {
	Type     string   `hcl:"type,label"`
	Name     string   `hcl:"name,label"`
	Compare  bool     `hcl:"compare"`
	Required bool     `hcl:"required"`
	Remain   hcl.Body `hcl:",remain"`
}

// ReadConfig reads the HCL config file and returns the parsed structure.
func ReadConfig(path string) (*HclConfig, error) {
	parser := hclparse.NewParser()

	hcfile, diag := parser.ParseHCLFile(path)
	if diag.HasErrors() {
		var errStrs []string
		for _, err := range diag.Errs() {
			errStrs = append(errStrs, err.Error())
		}
		combinedErr := strings.Join(errStrs, ": ")

		return nil, fmt.Errorf("read config %s: %s", path, combinedErr)
	}

	var config HclConfig
	diag = gohcl.DecodeBody(hcfile.Body, nil, &config)
	if diag.HasErrors() {
		var errStrs []string
		for _, err := range diag.Errs() {
			errStrs = append(errStrs, err.Error())
		}
		combinedErr := strings.Join(errStrs, ": ")

		return nil, fmt.Errorf("read config: decode body: %s", combinedErr)
	}

	return &config, nil
}
