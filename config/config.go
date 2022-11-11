package config

import (
	"fmt"
	"os"
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
	konfText string
)

var log = logging.New().
	WithPrefix("config")

// Loads a YAML configuration file from the given path and overrides
// it with environment variables. If the file cannot be loaded or
// parsed as YAML, an error is returned. Requires a default config of any kind,
// will try to serialize the configuration to outConfig if present (needs to
// be a pointer to a struct).
func Load(path string, onReload ...func(*koanf.Koanf)) (*koanf.Koanf, error) {
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
	konfTextBytes, _ := fp.ReadBytes()
	konfText = string(konfTextBytes)

	err := fp.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Error(err)
			return
		}

		log.Debug("configuration changed, reloading")
		konf = koanf.New(".")
		_ = konf.Load(fp, yaml.Parser())

		konfTextBytes, _ := fp.ReadBytes()
		konfText = string(konfTextBytes)

		for _, cb := range onReload {
			cb(konf)
		}
	})
	if err != nil {
		return nil, err
	}

	_ = konf.Load(env.ProviderWithValue("ORFS_", ".", formatEnv), nil)

	return konf, nil
}

// Save the given text in the configuration file on disk.
func Save(text string) error {
	return os.WriteFile(konfPath, []byte(text), 0o644)
}

// Returns full configuration file contents.
func Text() string {
	return konfText
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
