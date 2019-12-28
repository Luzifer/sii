package main

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var errUserConfigNotFound = errors.New("User config not found")

type configFile struct {
	GameDirectories    map[string]string `yaml:"game_directories"`
	ProfileDirectories map[string]string `yaml:"profile_directories"`
}

func (c *configFile) loadDefaults() error {
	for k, dv := range gamePaths {
		if c.GameDirectories[k] == "" {
			c.GameDirectories[k] = dv
		}
	}

	for k, dv := range profilePaths {
		if c.ProfileDirectories[k] == "" {
			c.ProfileDirectories[k] = dv
		}
	}

	return nil
}

func loadUserConfig(p string) (*configFile, error) {
	var c = &configFile{
		GameDirectories:    map[string]string{},
		ProfileDirectories: map[string]string{},
	}

	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return c, errUserConfigNotFound
		}
		return c, errors.Wrap(err, "Unable to stat user config")
	}

	f, err := os.Open(p)
	if err != nil {
		return c, errors.Wrap(err, "Unable to open user config")
	}
	defer f.Close()

	return c, errors.Wrap(yaml.NewDecoder(f).Decode(&c), "Unable to read user config")
}
