package meta

import (
	_ "embed"
	"strings"

	"gopkg.in/yaml.v3"
)

func init() {
	G.Name = "mymyapp"
	G.Version = parseVersion()
}

//go:embed meta.yml
var yamlFile string

func parseVersion() string {
	var meta struct {
		Version string `yaml:"version"`
	}

	if err := yaml.Unmarshal([]byte(yamlFile), &meta); err != nil {
		panic("failed to fetch version")
	}
	return strings.TrimSpace(meta.Version)
}
