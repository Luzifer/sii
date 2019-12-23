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

type commSaveDetails struct {
	CargoDamage     float32
	TruckWear       float32
	TrailerAttached bool
	TrailerWear     float32

	// FIXME: Add more details for profile
	// e.g. current job, money, xp
}

func commSaveDetailsFromUnit(unit *sii.Unit) (out commSaveDetails, err error) {
	var player *sii.Player

	for _, b := range unit.Entries {
		if v, ok := b.(*sii.Player); ok {
			player = v
			break
		}
	}

	if player == nil {
		return out, errors.New("Found no player object")
	}

	var (
		truck   *sii.Vehicle
		trailer *sii.Trailer
	)

	if v, ok := player.AssignedTruck.Resolve().(*sii.Vehicle); ok {
		truck = v
	}

	if v, ok := player.AssignedTrailer.Resolve().(*sii.Trailer); ok {
		trailer = v
	}

	if trailer != nil {
		for _, pb := range trailer.Accessories {
			var wear float32
			if v, ok := pb.Resolve().(*sii.VehicleAccessory); ok {
				wear = v.Wear
			}

			if v, ok := pb.Resolve().(*sii.VehicleWheelAccessory); ok {
				wear = v.Wear
			}

			if wear > out.TrailerWear {
				out.TrailerWear = wear
			}
		}

		out.CargoDamage = trailer.CargoDamage
	}

	for _, pb := range truck.Accessories {
		var wear float32
		if v, ok := pb.Resolve().(*sii.VehicleAccessory); ok {
			wear = v.Wear
		}

		if v, ok := pb.Resolve().(*sii.VehicleWheelAccessory); ok {
			wear = v.Wear
		}

		if wear > out.TruckWear {
			out.TruckWear = wear
		}
	}

	return out, nil
}

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
