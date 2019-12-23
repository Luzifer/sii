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

type commProfileInfo struct {
	CompanyName  string
	ProfileName  string
	CreationTime time.Time
	SaveTime     time.Time
}

func commProfileInfoFromUserProfile(p *sii.UserProfile) commProfileInfo {
	var out commProfileInfo

	out.CompanyName = p.CompanyName
	out.ProfileName = p.ProfileName
	out.CreationTime = time.Unix(p.CreationTime, 0)
	out.SaveTime = time.Unix(p.SaveTime, 0)

	return out
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
