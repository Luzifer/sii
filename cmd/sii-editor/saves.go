package main

import (
	"os"
	"path"
	"strconv"
	"strings"
	"time"

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
	GameTime int64 `json:"game_time"`

	ExperiencePoints int64 `json:"experience_points"`
	Money            int64 `json:"money"`

	CargoDamage     float32 `json:"cargo_damage"`
	TruckWear       float32 `json:"truck_wear"`
	TrailerAttached bool    `json:"trailer_attached"`
	TrailerWear     float32 `json:"trailer_wear"`

	AttachedTrailer string            `json:"attached_trailer"`
	OwnedTrailers   map[string]string `json:"owned_trailers"`

	CurrentJob *commSaveJob `json:"current_job"`
}

type commSaveJob struct {
	OriginReference  string  `json:"origin_reference"`
	TargetReference  string  `json:"target_reference"`
	CargoReference   string  `json:"cargo_reference"`
	TrailerReference string  `json:"trailer_reference"`
	Distance         int64   `json:"distance"`
	Urgency          *int64  `json:"urgency,omitempty"`
	Weight           float32 `json:"weight,omitempty"`
	Expires          int64   `json:"expires,omitempty"`
}

func commSaveDetailsFromUnit(unit *sii.Unit) (out commSaveDetails, err error) {
	var economy *sii.Economy

	for _, b := range unit.BlocksByClass("economy") {
		if v, ok := b.(*sii.Economy); ok {
			economy = v
		}
	}

	if economy == nil {
		return out, errors.New("Found no economy object")
	}

	var (
		bank    *sii.Bank
		job     *sii.PlayerJob
		player  *sii.Player
		truck   *sii.Vehicle
		trailer *sii.Trailer
	)

	bank = economy.Bank.Resolve().(*sii.Bank)
	player = economy.Player.Resolve().(*sii.Player)

	out.ExperiencePoints = economy.ExperiencePoints
	out.GameTime = economy.GameTime
	out.Money = bank.MoneyAccount
	out.TrailerAttached = player.AssignedTrailerConnected

	if v, ok := player.CurrentJob.Resolve().(*sii.PlayerJob); ok {
		job = v
	}

	if v, ok := player.AssignedTruck.Resolve().(*sii.Vehicle); ok {
		truck = v
	}

	if v, ok := player.AssignedTrailer.Resolve().(*sii.Trailer); ok {
		out.AttachedTrailer = player.AssignedTrailer.Target
		trailer = v
	}

	if truck != nil {
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
	}

	if len(player.Trailers) > 0 {
		out.OwnedTrailers = map[string]string{}
		for _, tp := range player.Trailers {
			if t, ok := tp.Resolve().(*sii.Trailer); ok {
				out.OwnedTrailers[t.Name()] = t.CleanedLicensePlate()
			}
		}
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

	if job != nil {
		out.CurrentJob = &commSaveJob{
			OriginReference:  job.SourceCompany.Target,
			TargetReference:  job.TargetCompany.Target,
			CargoReference:   job.Cargo.Target,
			TrailerReference: job.CompanyTrailer.Target,
			Distance:         job.PlannedDistanceKM,
			Urgency:          job.Urgency,
		}
	}

	return out, nil
}

func addJobToGame(game *sii.Unit, job commSaveJob, routePartID int) error {
	// Get company
	company := game.BlockByName(job.OriginReference).(*sii.Company)
	// Get cargo
	cargo := baseGameUnit.BlockByName(job.CargoReference).(*sii.CargoData)

	// Get trailer / truck from other jobs
	var cTruck, cTV, cTD string

	for _, jp := range company.JobOffer {
		j := jp.Resolve().(*sii.JobOfferData)
		if j.CompanyTruck.IsNull() || j.TrailerVariant.IsNull() || j.TrailerDefinition.IsNull() {
			continue
		}
		cTruck, cTV, cTD = j.CompanyTruck.Target, j.TrailerVariant.Target, j.TrailerDefinition.Target
		break
	}

	if cTruck == "" || cTV == "" || cTD == "" {
		// The company did not have any valid job offers to steal from, lets search globally
		for _, jb := range game.BlocksByClass("job_offer_data") {
			j := jb.(*sii.JobOfferData)
			if j.CompanyTruck.IsNull() || j.TrailerVariant.IsNull() || j.TrailerDefinition.IsNull() {
				continue
			}
			cTruck, cTV, cTD = j.CompanyTruck.Target, j.TrailerVariant.Target, j.TrailerDefinition.Target
			break
		}
	}

	if cTruck == "" || cTV == "" || cTD == "" {
		// Even in global jobs there was no proper job, this savegame looks broken
		return errors.New("Was not able to find suitable trailer definition")
	}

	cargoCount := int64(job.Weight / cargo.Mass)
	if cargoCount < 1 {
		// Ensure cargo is transported (can happen if only small weight is
		// requested and one unit weighs more than requested cargo)
		cargoCount = 1
	}

	jobID := strings.Join([]string{
		"_nameless",
		strconv.FormatInt(time.Now().Unix(), 16),
		strconv.FormatInt(int64(routePartID), 16),
	}, ".")
	exTime := game.BlocksByClass("economy")[0].(*sii.Economy).GameTime + 300 // 300min = 5h
	j := &sii.JobOfferData{
		// User requested job data
		Target:             strings.TrimPrefix(job.TargetReference, "company.volatile."),
		ExpirationTime:     &exTime,
		Urgency:            job.Urgency,
		Cargo:              sii.Ptr{Target: job.CargoReference},
		UnitsCount:         cargoCount,
		ShortestDistanceKM: job.Distance,

		// Some static data
		FillRatio: 1, // Dunno but other jobs have it at 1, so keep for now

		// Dunno where this data comes from, it works without it and gets own ones
		TrailerPlace: []sii.Placement{},

		// Too lazy to implement, just steal it too
		CompanyTruck:      sii.Ptr{Target: cTruck},
		TrailerVariant:    sii.Ptr{Target: cTV},
		TrailerDefinition: sii.Ptr{Target: cTD},
	}
	j.Init("", jobID)

	// Add the new job to the game
	game.Entries = append(game.Entries, j)
	company.JobOffer = append([]sii.Ptr{{Target: j.Name()}}, company.JobOffer...)

	return nil
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
