package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	APIURL string `toml:"api_url"`
	APIKey string `toml:"api_key"`
}

func Path() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "hizal", "config.toml")
}

func Load() (*Config, error) {
	path := Path()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{}, nil
	}
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path := Path()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

func (c *Config) IsConfigured() bool {
	return c.APIURL != "" && c.APIKey != ""
}
