package shared

import (
	"bytes"
	"encoding/json"
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
	Name   string `json:"name"`
	Unique bool   `json:"unique"`
}

type MinerPropertyContent struct {
	Format string `json:"format"`
	Value  string `json:"value"`
}

type MinerProperty struct {
	Type    string               `json:"type"`
	Label   MinerPropertyLabel   `json:"label"`
	Content MinerPropertyContent `json:"content"`
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

type MinerOutline struct {
	Resource MinerResource `json:"resource"`
	Diary    MinerDiary    `json:"diary"`
}

type MinerDiary struct {
	Hash string `json:"hash"`
	Logs struct {
		Prev string `json:"prev"`
		Curr string `json:"curr"`
	} `json:"logs"`
}

func NewMinerDiary(hash, currLog, prevLog string) MinerDiary {
	return MinerDiary{
		Hash: hash,
		Logs: struct {
			Prev string `json:"prev"`
			Curr string `json:"curr"`
		}{
			Prev: prevLog,
			Curr: currLog,
		},
	}
}

type MinerResource struct {
	Identifier string          `json:"identifier"`
	Alias      string          `json:"alias"`
	LogType    string          `json:"logType"`
	Properties []MinerProperty `json:"properties"`
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

func (m *MinerResource) RenderMarkdown() (string, error) {
	var sb strings.Builder

	if _, err := sb.WriteString(fmt.Sprintf("# Identifier\n**%s**\n\n", m.Identifier)); err != nil {
		return "", fmt.Errorf("MinerResource RenderMarkdown(): %w", err)
	}
	if m.Alias != "" {
		if _, err := sb.WriteString(fmt.Sprintf("## Alias\n%s\n\n", m.Alias)); err != nil {
			return "", fmt.Errorf("MinerResource RenderMarkdown(): %w", err)
		}
	}

	for i, prop := range m.Properties {
		if i == 0 {
			if _, err := sb.WriteString("## Properties\n\n"); err != nil {
				return "", fmt.Errorf("MinerResource RenderMarkdown(): %w", err)
			}
		}

		if _, err := sb.WriteString(fmt.Sprintf("### %s\n", prop.Type)); err != nil {
			return "", fmt.Errorf("MinerResource RenderMarkdown(): %w", err)
		}
		if _, err := sb.WriteString(fmt.Sprintf("- **Label:** %s\n", prop.Label.Name)); err != nil {
			return "", fmt.Errorf("MinerResource RenderMarkdown(): %w", err)
		}
		if _, err := sb.WriteString("- **Content:**\n"); err != nil {
			return "", fmt.Errorf("MinerResource RenderMarkdown(): %w", err)
		}
		if prop.Content.Format == FormatJson {
			prettyJson := bytes.Buffer{}
			if err := json.Indent(&prettyJson, []byte(prop.Content.Value), "", "  "); err != nil {
				return "", fmt.Errorf("MinerResource RenderMarkdown(): %w", err)
			}
			input := indentString(prettyJson.String(), "  ")
			if _, err := sb.WriteString(fmt.Sprintf("  ```json\n%s\n  ```\n", input)); err != nil {
				return "", fmt.Errorf("MinerResource RenderMarkdown(): %w", err)
			}
		} else if prop.Content.Format == FormatText {
			input := indentString(prop.Content.Value, "  ")
			if _, err := sb.WriteString(fmt.Sprintf("  ```\n%s\n  ```\n", input)); err != nil {
				return "", fmt.Errorf("MinerResource RenderMarkdown(): %w", err)
			}
		}
	}

	return sb.String(), nil
}

func indentString(input string, indent string) string {
	// Split the input string by newline
	lines := strings.Split(input, "\n")

	// Prepend each line with the indent string
	for i, line := range lines {
		lines[i] = indent + line
	}

	// Join the lines back together with newline characters
	return strings.Join(lines, "\n")
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
