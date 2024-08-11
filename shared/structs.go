package shared

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type MinerConfigEquipment struct {
	Type       string
	Name       string
	Attributes map[string]string
}

type MinerConfig struct {
	Auth       map[string]string
	Equipments []MinerConfigEquipment
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

// FormatContentValue set input data to desired property content value based on the content format.
func (m *MinerProperty) FormatContentValue(data any) error {
	switch m.Content.Format {
	case FormatJson:
		marshaledData, err := JsonMarshal(data)
		if err != nil {
			return fmt.Errorf("MinerProperty format: %w", err)
		}

		normalizedJson, err := JsonNormalize(string(marshaledData))
		if err != nil {
			return fmt.Errorf("MinerProperty format: %w", err)
		}
		m.Content.Value = string(normalizedJson)
	case FormatText:
		m.Content.Value = fmt.Sprintf("%s", data)
	default:
		return fmt.Errorf("MinerProperty format: unknown format: %s", m.Content.Format)
	}

	return nil
}

type MinerResource struct {
	Identifier string
    Alias string
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
	Name          string            `hcl:"name,label"`
	Group         string            `hcl:"group,label"`
	Authenticator map[string]string `hcl:"authenticator,attr"`
	Diaries       []PlugDiary       `hcl:"diary,block"`
	Equipments    []PlugEquipment   `hcl:"equipment,block"`
}

func (p Plug) GenMinerConfig() MinerConfig {
	equipments := []MinerConfigEquipment{}
	for _, equipment := range p.Equipments {
		equipments = append(equipments, MinerConfigEquipment{
			Type:       equipment.Type,
			Name:       equipment.Name,
			Attributes: equipment.Attributes,
		})
	}

	return MinerConfig{
		Auth:       p.Authenticator,
		Equipments: equipments,
	}
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
