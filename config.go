package main

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type httpsConfig struct {
	Key  string `yaml:"key"`
	Cert string `yaml:"cert"`
}

// FileConfig struct represents config for a single file.
type FileConfig struct {
	Cache bool     `yaml:"cache"`
	Gzip  bool     `yaml:"gzip"`
	Push  []string `yaml:"push"`
}

// Config struct represents runtime config for statik.
type Config struct {
	Listen string                 `yaml:"listen"`
	Files  map[string]*FileConfig `yaml:"files"`
	HTTPS  *httpsConfig           `yaml:"https,omitempty"`
	Root   string
}

// NewConfig func creates new config.
func NewConfig() (*Config, error) {
	c := &Config{}

	pwd, err := getPwd()
	if err != nil {
		return nil, fmt.Errorf("cannot read pwd: %q", err)
	}
	c.Root = pwd

	return c, nil
}

// NewConfigFromYaml func creates config from given yaml bytes.
func NewConfigFromYaml(b []byte) (*Config, error) {
	c, err := NewConfig()

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config: %q", err)
	}

	return c, nil
}

// IsHTTPS func returns whether https is on.
func (c *Config) IsHTTPS() bool { return c.HTTPS != nil }

// GetConfigForPath func returns FileConfig for given path.
func (c *Config) GetConfigForPath(path string) *FileConfig {
	p, ok := c.Files[path]

	if ok {
		return p
	}

	// TODO: look wildcard configs.

	return &FileConfig{}
}

func getPwd() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	return path.Dir(ex), nil
}
