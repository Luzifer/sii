package main

import (
	"os"
	"path"

	"github.com/Luzifer/sii"
	"github.com/pkg/errors"
)

const (
	fixAll         = "all"
	fixCargoOnly   = "cargo"
	fixTrailerOnly = "trailer"
	fixTruckOnly   = "truck"
)

func fixPlayerTruck(unit *sii.Unit, fixType string) error {
	// In call cases we need the player as the starting point
	var player *sii.Player

	for _, b := range unit.Entries {
		if v, ok := b.(*sii.Player); ok {
			player = v
			break
		}
	}

	if player == nil {
		return errors.New("Found no player object")
	}

	// Fix truck
	if fixType == fixAll || fixType == fixTruckOnly {
		truck := player.AssignedTruck.Resolve().(*sii.Vehicle)
		for _, pb := range truck.Accessories {
			if v, ok := pb.Resolve().(*sii.VehicleAccessory); ok {
				v.Wear = 0
			}

			if v, ok := pb.Resolve().(*sii.VehicleWheelAccessory); ok {
				v.Wear = 0
			}
		}
	}

	// Fix trailer
	if (fixType == fixAll || fixType == fixTrailerOnly) && player.AssignedTrailer.Resolve() != nil {
		trailer := player.AssignedTrailer.Resolve().(*sii.Trailer)
		for _, pb := range trailer.Accessories {
			if v, ok := pb.Resolve().(*sii.VehicleAccessory); ok {
				v.Wear = 0
			}

			if v, ok := pb.Resolve().(*sii.VehicleWheelAccessory); ok {
				v.Wear = 0
			}
		}
	}

	// Fix cargo
	if (fixType == fixAll || fixType == fixCargoOnly) && player.AssignedTrailer.Resolve() != nil {
		trailer := player.AssignedTrailer.Resolve().(*sii.Trailer)
		trailer.CargoDamage = 0
	}

	return nil
}

// loadSave reads game- and info- unit and returns them unmodified
func loadSave(profile, save string) (*sii.Unit, *sii.SaveContainer, error) {
	var (
		saveInfoPath = getSaveFilePath(profile, save, "info.sii")
		saveUnitPath = getSaveFilePath(profile, save, "game.sii")
	)

	infoUnit, err := sii.ReadUnitFile(saveInfoPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to load save-info")
	}

	var info *sii.SaveContainer
	for _, b := range infoUnit.Entries {
		if v, ok := b.(*sii.SaveContainer); ok {
			info = v
		}
	}

	gameUnit, err := sii.ReadUnitFile(saveUnitPath)
	if err != nil {
		return nil, info, errors.Wrap(err, "Unable to load save-unit")
	}

	return gameUnit, info, nil
}

// storeSave writes game- and info- unit without checking for existance!
func storeSave(profile, save string, unit *sii.Unit, info *sii.SaveContainer) error {
	var (
		saveInfoPath = getSaveFilePath(profile, save, "info.sii")
		saveUnitPath = getSaveFilePath(profile, save, "game.sii")
	)

	if err := os.MkdirAll(path.Dir(saveInfoPath), 0700); err != nil {
		return errors.Wrap(err, "Unable to create save-dir")
	}

	if err := sii.WriteUnitFile(saveInfoPath, &sii.Unit{Entries: []sii.Block{info}}); err != nil {
		return errors.Wrap(err, "Unable to write info unit")
	}

	if err := sii.WriteUnitFile(saveUnitPath, unit); err != nil {
		return errors.Wrap(err, "Unable to write game unit")
	}

	return nil
}
