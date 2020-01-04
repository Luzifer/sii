package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/scs-extract/scs"
	"github.com/Luzifer/sii"
)

var baseDataFiles = regexp.MustCompile(`^def/(?:cargo|city|company)/[^/]+.sii$`)

func readBaseData() (*sii.Unit, error) {
	var unitData = new(bytes.Buffer)
	// Open a plain unit for parsing
	unitData.WriteString("SiiNunit\n{\n")

	// Collect all available units from game files
	entries, err := ioutil.ReadDir(cfg.GameDir)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to list game-directory")
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".scs") {
			// We don't care for anthing than SCS# files
			continue
		}

		fPath := path.Join(cfg.GameDir, entry.Name())

		stat, err := os.Stat(fPath)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to stat SCS# archive")
		}

		scsFile, err := os.Open(fPath)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to open SCS# archive")
		}
		defer scsFile.Close()

		r, err := scs.NewReader(scsFile, stat.Size())
		if err != nil {
			// There are some files being unreadable, that's okay
			log.WithField("file", entry.Name()).Debug("Found unreadable SCS archive")
			continue
		}

		for _, f := range r.Files {
			if !baseDataFiles.MatchString(f.Name) {
				// We don't care for most of the files, just mentioned definitions
				continue
			}

			fr, err := f.Open()
			if err != nil {
				return nil, errors.Wrap(err, "Unable to open file from SCS archive")
			}

			if _, err = io.Copy(unitData, fr); err != nil {
				return nil, errors.Wrap(err, "Unable to read file from SCS archive")
			}

			unitData.WriteString("\n") // Ensure CR after each block

			fr.Close()
		}
	}

	// Close unit
	unitData.WriteString("}")

	// Read-in constructed unit file
	return sii.ParseSIIPlainFile(unitData)
}
