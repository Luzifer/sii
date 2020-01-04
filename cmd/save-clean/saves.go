package main

import (
	"path"

	"github.com/pkg/errors"

	"github.com/Luzifer/sii"
)

// loadSave reads game- and info- unit and returns them unmodified
func loadSave(gameSIIPath string) (*sii.Unit, error) {
	gameUnit, err := sii.ReadUnitFile(gameSIIPath)
	return gameUnit, errors.Wrap(err, "Unable to load save-unit")
}

// storeSave writes game- and info- unit without checking for existance!
func storeSave(gameSIIPath string, unit *sii.Unit) error {
	var saveUnitPath = path.Join(path.Dir(gameSIIPath), "game.cleaned.sii")

	if err := sii.WriteUnitFile(saveUnitPath, unit); err != nil {
		return errors.Wrap(err, "Unable to write game unit")
	}

	return nil
}
