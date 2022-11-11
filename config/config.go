package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/openrfsense/common/logging"
)

type Collector struct {
	Port int `yaml:"port"`
}

type Location struct {
	Name      string  `yaml:"name"`
	Elevation float64 `yaml:"elevation"`
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
}

type Node struct {
	Port int `yaml:"port"`
}

type NATS struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type NodeConfig struct {
	Collector `yaml:"collector"`
	Location  `yaml:"location"`
	Node      `yaml:"node"`
	NATS      `yaml:"nats"`
}

var defaultConfig = NodeConfig{
	Collector: Collector{
		Port: 2022,
	},
	Node: Node{
		Port: 9090,
	},
	NATS: NATS{
		Port: 0,
	},
}

var (
	konf     *koanf.Koanf
	konfPath string
)

var log = logging.New().
	WithPrefix("config")

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

	fp := file.Provider(path)
	if err := konf.Load(fp, yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading configuration file: %v (%T)", err, err)
	}

	konfPath = path

	err := fp.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Error(err)
			return
		}

		log.Debug("configuration changed, reloading")
		konf = koanf.New(".")
		_ = konf.Load(fp, yaml.Parser())
	})
	if err != nil {
		return nil, err
	}

	_ = konf.Load(env.ProviderWithValue("ORFS_", ".", formatEnv), nil)

	return konf, nil
}

// Returns a copy of the configuration file path in use.
func Path() string {
	return konfPath
}

// Formats environment variables: ORFS_SECTION_SUBSECTION_KEY becomes
// (as a path) section.subsection.key
func formatEnv(s string, v string) (string, interface{}) {
	rawPath := strings.ToLower(strings.TrimPrefix(s, "ORFS_"))
	key := strings.Replace(rawPath, "_", ".", -1)

	if strings.Contains(v, " ") {
		return key, strings.Split(v, " ")
	}

	return key, v
}
