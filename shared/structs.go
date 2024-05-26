package shared

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type MinerConfig struct {
	Path string
}

type MinerPropertyLabel struct {
	Name   string
	Unique bool
}

type MinerPropertyContent struct {
	Format string
	Value  string
}

type MinerProperty struct {
	Type    string
	Label   MinerPropertyLabel
	Content MinerPropertyContent
}

type MinerResource struct {
	Identifier string
	Properties []MinerProperty
}

func (m *MinerResource) Sort() {
    sort.SliceStable(m.Properties, func(i, j int) bool {
        if m.Properties[i].Type == m.Properties[j].Type {
            if m.Properties[i].Label.Name == m.Properties[j].Label.Name {
                return m.Properties[i].Content.Value < m.Properties[j].Content.Value
            }
            return m.Properties[i].Label.Name < m.Properties[j].Label.Name
        }
        return m.Properties[i].Type < m.Properties[j].Type
    })
}

type MinerResources []MinerResource


// HCL config structure
type HclConfig struct {
	Plugs []Plug `hcl:"plug,block"`
}

type Plug struct {
	Name       string     `hcl:"name,label"`
	Group      string     `hcl:"group,label"`
	Profile    string     `hcl:"profile"`
	Properties []Property `hcl:"property,block"`
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
