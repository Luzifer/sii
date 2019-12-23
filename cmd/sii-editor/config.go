package main

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var errUserConfigNotFound = errors.New("User config not found")

type configFile struct {
	GameDirectory    string `yaml:"game_directory"`
	ProfileDirectory string `yaml:"profile_directory"`
}

func (c *configFile) loadDefaults() error {
	var err error

	if c.GameDirectory == "" {
		if c.GameDirectory, err = findGamePath(); err != nil {
			return errors.Wrap(err, "Unable to find game directory")
		}
	}

	if c.ProfileDirectory == "" {
		if c.ProfileDirectory, err = findProfilePath(); err != nil {
			return errors.Wrap(err, "Unable to find profile directory")
		}
	}

	return nil
}

func loadUserConfig(p string) (*configFile, error) {
	var c = &configFile{}

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
