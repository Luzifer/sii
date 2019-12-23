package main

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

var errPathNotFound = errors.New("Could not find path")

func findProfilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "Unable to get users homedir")
	}

	for _, hintPath := range profilePaths {
		hintPath = strings.ReplaceAll(hintPath, "~", homedir)
		if s, err := os.Stat(hintPath); err == nil {
			if s.IsDir() {
				return hintPath, nil
			}
		}
	}

	return "", errPathNotFound
}

func findGamePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "Unable to get users homedir")
	}

	for _, hintPath := range gamePaths {
		hintPath = strings.ReplaceAll(hintPath, "~", homedir)
		if _, err := os.Stat(path.Join(hintPath, "def.scs")); err == nil {
			return hintPath, nil
		}
	}

	return "", errPathNotFound
}

func getProfilePath(profile string) string {
	return path.Join(userConfig.ProfileDirectory, profile)
}

func getProfileInfoPath(profile string) string {
	return path.Join(getProfilePath(profile), "profile.sii")
}

func getSavePath(profile, save string) string {
	return path.Join(getProfilePath(profile), "save", save)
}

func getSaveFilePath(profile, save, file string) string {
	return path.Join(getSavePath(profile, save), file)
}
