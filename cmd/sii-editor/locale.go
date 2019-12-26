package main

import (
	"bytes"
	"io"
	"os"
	"path"

	"github.com/Luzifer/scs-extract/scs"
	"github.com/Luzifer/sii"
	"github.com/Luzifer/sii/t3nk"
	"github.com/pkg/errors"
)

func getLocale(locale string) (*sii.LocalizationDB, error) {
	fPath := path.Join(userConfig.GameDirectory, "locale.scs")

	stat, err := os.Stat(fPath)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to stat locale.scs archive")
	}

	scsFile, err := os.Open(fPath)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open locale.scs archive")
	}
	defer scsFile.Close()

	r, err := scs.NewReader(scsFile, stat.Size())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open locale.scs archive")
	}

	var codedLocale = new(bytes.Buffer)
	for _, f := range r.Files {
		if f.Name == path.Join("locale", locale, "local.sii") {
			lf, err := f.Open()
			if err != nil {
				return nil, errors.Wrap(err, "Unable to read local.sii file from archive")
			}
			defer lf.Close()

			// Using the reader directly somehow broke mid-file, pre-buffering works fine
			if _, err := io.Copy(codedLocale, lf); err != nil {
				return nil, errors.Wrap(err, "Unable to copy local.sii file from archive")
			}
			break
		}
	}

	if codedLocale.Len() == 0 {
		return nil, errors.New("Found no locale information")
	}

	localeReader, err := t3nk.Decode(codedLocale)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to decode locale.sii")
	}

	unit, err := sii.ParseSIIPlainFile(localeReader)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read unit from locale.sii")
	}

	for _, b := range unit.BlocksByClass("localization_db") {
		if v, ok := b.(*sii.LocalizationDB); ok {
			return v, nil
		}
	}

	return nil, errors.New("No localization db found in locale")
}
