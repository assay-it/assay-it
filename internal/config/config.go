package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultConfigFile = ".assay-it.json"
)

// Suits config
type Config struct {
	Runner string   `json:"runner"`
	Suites []string `json:"suites"`
}

func NewFromPkg(pkg string) (*Config, error) {
	if pkg == "" {
		return parseConfigFileJSON(defaultConfigFile)
	}

	return parseConfigFileJSON(filepath.Join(pkg, defaultConfigFile))
}

func parseConfigFileJSON(file string) (*Config, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	if len(config.Suites) == 0 {
		return nil, fmt.Errorf("suites are not defined at %s", file)
	}

	if config.Runner == "" {
		runner := filepath.Dir(config.Suites[0])
		config.Runner = filepath.Join(runner, "assay-it")
	}

	return &config, nil
}
