package main

import (
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	pathETS2 = "ets2"
	pathATS  = "ats"
)

var errPathNotFound = errors.New("Could not find path")

func expandHomedir(dir string) string {
	s, err := homedir.Expand(dir)
	if err != nil {
		log.WithError(err).Error("Unable to expand home path")
		return dir
	}
	return s
}

func getGamePath() string {
	return expandHomedir(userConfig.GameDirectories[cfg.Game])
}

func getProfilesPath() string {
	return expandHomedir(userConfig.ProfileDirectories[cfg.Game])
}

func getProfilePath(profile string) string {
	return path.Join(getProfilesPath(), profile)
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
