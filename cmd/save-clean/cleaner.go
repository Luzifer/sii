package main

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/sii"
)

func cleanSaveGame(baseGameUnit, saveGame *sii.Unit) error {
	log.Info("Cleaning companies...")
	if err := cleanCompanies(baseGameUnit, saveGame); err != nil {
		return errors.Wrap(err, "Unable to cleanup companies")
	}

	log.Info("Cleaning garages...")
	if err := cleanGarages(baseGameUnit, saveGame); err != nil {
		return errors.Wrap(err, "Unable to cleanup garages")
	}

	return nil
}
