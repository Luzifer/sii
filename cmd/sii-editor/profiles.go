package main

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Luzifer/sii"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type commSaveInfo struct {
	Name     string    `json:"name"`
	GameTime int64     `json:"game_time"`
	SaveTime time.Time `json:"file_type"`
}

func commSaveInfoFromSaveContainer(c *sii.SaveContainer) commSaveInfo {
	var out commSaveInfo

	out.Name = c.SaveName
	out.GameTime = c.Time
	out.SaveTime = time.Unix(c.FileTime, 0)

	return out
}

type commProfileInfo struct {
	CompanyName  string    `json:"company_name"`
	ProfileName  string    `json:"profile_name"`
	CreationTime time.Time `json:"creation_time"`
	SaveTime     time.Time `json:"save_time"`
}

func commProfileInfoFromUserProfile(p *sii.UserProfile) commProfileInfo {
	var out commProfileInfo

	out.CompanyName = p.CompanyName
	out.ProfileName = p.ProfileName
	out.CreationTime = time.Unix(p.CreationTime, 0)
	out.SaveTime = time.Unix(p.SaveTime, 0)

	return out
}

func listSaves(profile string) (map[string]commSaveInfo, error) {
	entries, err := ioutil.ReadDir(path.Join(userConfig.ProfileDirectory, profile, "save"))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to list saves")
	}

	var out = map[string]commSaveInfo{}

	for _, entry := range entries {
		if !entry.IsDir() {
			// There shouldn't be files but whatever
			continue
		}

		var sFile = path.Join(userConfig.ProfileDirectory, profile, "save", entry.Name(), "info.sii")
		if _, err := os.Stat(sFile); err != nil {
			// That directory contains no profile.sii - Weird but okay
			log.WithFields(log.Fields{
				"profile": profile,
				"save":    entry.Name(),
			}).Debug("Found save directory without info.sii")
			continue
		}

		unit, err := sii.ReadUnitFile(sFile)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to read unit for save %q", entry.Name())
		}

		for _, b := range unit.Entries {
			if v, ok := b.(*sii.SaveContainer); ok {
				out[entry.Name()] = commSaveInfoFromSaveContainer(v)
				break
			}
		}
	}

	return out, nil
}

func listProfiles() (map[string]commProfileInfo, error) {
	entries, err := ioutil.ReadDir(userConfig.ProfileDirectory)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to list profiles")
	}

	var out = map[string]commProfileInfo{}

	for _, entry := range entries {
		if !entry.IsDir() {
			// There shouldn't be files but whatever
			continue
		}

		var pFile = path.Join(userConfig.ProfileDirectory, entry.Name(), "profile.sii")
		if _, err := os.Stat(pFile); err != nil {
			// That directory contains no profile.sii - Weird but okay
			log.WithField("profile", entry.Name()).Debug("Found profile directory without profile.sii")
			continue
		}

		unit, err := sii.ReadUnitFile(pFile)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to read unit for profile %q", entry.Name())
		}

		for _, b := range unit.Entries {
			if v, ok := b.(*sii.UserProfile); ok {
				out[entry.Name()] = commProfileInfoFromUserProfile(v)
				break
			}
		}
	}

	return out, nil
}
