package main

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// ConfigLoader ...
type ConfigLoader interface {
	Load(string) (*Config, error)
}

type configLoaderImpl struct {
}

// NewConfigLoader ...
func NewConfigLoader() ConfigLoader {
	return &configLoaderImpl{}
}

func (loader *configLoaderImpl) Load(path string) (*Config, error) {
	fp := filepath.Clean(path)
	bytes, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
