package config

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Module string   `json:"module"`
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

	if config.Module == "" {
		path, err := filepath.Abs(file)
		if err != nil {
			return nil, err
		}

		config.Module, err = discoverModule(filepath.Dir(path))
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func discoverModule(path string) (string, error) {
	if path == string(filepath.Separator) {
		return "", nil
	}

	mod := filepath.Join(path, "go.mod")
	if _, err := os.Stat(mod); errors.Is(err, os.ErrNotExist) {
		return discoverModule(filepath.Dir(path))
	}

	spec, err := os.ReadFile(mod)
	if err != nil {
		return "", err
	}

	seq := bytes.Split(spec, []byte("\n"))

	if !bytes.HasPrefix(seq[0], []byte("module ")) {
		return "", fmt.Errorf("invalid go.mod %s", mod)
	}

	return string(bytes.TrimPrefix(seq[0], []byte("module "))), nil
}
