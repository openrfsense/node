package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

type Node struct {
	Port int `koanf:"port"`
}

type Backend struct {
	Port  int               `koanf:"port"`
	Users map[string]string `koanf:"users"`
}

type NATS struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type NodeConfig struct {
	Node    `koanf:"node"`
	Backend `koanf:"backend"`
	NATS    `koanf:"nats"`
}

var defaultConfig = NodeConfig{
	Node: Node{
		Port: 8080,
	},
	Backend: Backend{
		Port: 8080,
	},
	NATS: NATS{
		Port: 0,
	},
}

var konf *koanf.Koanf

// Loads a YAML configuration file from the given path and overrides
// it with environment variables. If the file cannot be loaded or
// parsed as YAML, an error is returned. Requires a default config of any kind,
// will try to serialize the configuration to outConfig if present (needs to
// be a pointer to a struct).
func Load(path string) (*koanf.Koanf, error) {
	konf = koanf.New(".")

	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("configuration file path cannot be empty")
	}

	_ = konf.Load(structs.Provider(defaultConfig, "yaml"), nil)

	if err := konf.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading configuration file: %v (%T)", err, err)
	}

	_ = konf.Load(env.Provider("ORFS_", ".", formatEnv), nil)

	return konf, nil
}

// Formats environment variables: ORFS_SECTION_SUBSECTION_KEY becomes
// (as a path) section.subsection.key
func formatEnv(s string) string {
	rawPath := strings.ToLower(strings.TrimPrefix(s, "ORFS_"))
	return strings.Replace(rawPath, "_", ".", -1)
}
