package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Luzifer/go_helpers/v2/str"
	"github.com/Luzifer/sii"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const (
	storeSaveFolder = "sii_editor"
	storeSaveName   = "SII Editor"
)

func init() {
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}", handleGetSaveInfo).Methods(http.MethodGet)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/economy", handleUpdateEconomyInfo).Methods(http.MethodPut)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/fix", handleFixPlayer).Methods(http.MethodPut)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/jobs", handleListJobs).Methods(http.MethodGet)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/jobs", handleAddJob).Methods(http.MethodPost)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/set-trailer", handleSetTrailer).Methods(http.MethodPut)
}

func handleAddJob(w http.ResponseWriter, r *http.Request) {
	var (
		job  commSaveJob
		vars = mux.Vars(r)
	)

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		apiGenericError(w, http.StatusBadRequest, errors.Wrap(err, "Unable to decode input data"))
		return
	}

	game, info, err := loadSave(vars["profileID"], vars["saveFolder"])
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to load save"))
		return
	}

	info.SaveName = storeSaveName
	info.FileTime = time.Now().Unix()

	// Set urgency if it isn't
	if job.Urgency == nil {
		u := int64(0)
		job.Urgency = &u
	}

	if job.Weight == 0 {
		job.Weight = 10000 // 10 tons as a default
	}

	if job.Weight < 100 {
		// User clearly did't want 54kg but 54 tons! (If not: screw em)
		job.Weight *= 1000
	}

	if job.Distance == 0 {
		job.Distance = 100 // If the user did not provide distance use 100km as a default
	}

	// Get company
	company := game.BlockByName(job.OriginReference).(*sii.Company)
	// Get cargo
	cargo := baseGameUnit.BlockByName(job.CargoReference).(*sii.CargoData)

	// Get trailer / truck from other jobs
	var (
		cTruck, cTV, cTD string
		cTP              []sii.Placement
	)

	for _, jp := range company.JobOffer {
		j := jp.Resolve().(*sii.JobOfferData)
		if j.CompanyTruck.Target == "null" || j.TrailerVariant.Target == "null" || j.TrailerDefinition.Target == "null" || len(j.TrailerPlace) < 1 {
			continue
		}
		cTruck, cTV, cTD = j.CompanyTruck.Target, j.TrailerVariant.Target, j.TrailerDefinition.Target
		cTP = j.TrailerPlace
		break
	}

	if cTP == nil {
		// The company did not have any valid job offers to steal from, lets search globally
		for _, jb := range game.BlocksByClass("job_offer_data") {
			j := jb.(*sii.JobOfferData)
			if j.CompanyTruck.Target == "null" || j.TrailerVariant.Target == "null" || j.TrailerDefinition.Target == "null" || len(j.TrailerPlace) < 1 {
				continue
			}
			cTruck, cTV, cTD = j.CompanyTruck.Target, j.TrailerVariant.Target, j.TrailerDefinition.Target
			cTP = j.TrailerPlace
			break
		}
	}

	jobID := "_nameless." + strconv.FormatInt(time.Now().Unix(), 16)
	exTime := game.BlocksByClass("economy")[0].(*sii.Economy).GameTime + 300 // 300min = 5h
	j := &sii.JobOfferData{
		// User requested job data
		Target:             strings.TrimPrefix(job.TargetReference, "company.volatile."),
		ExpirationTime:     &exTime,
		Urgency:            job.Urgency,
		Cargo:              sii.Ptr{Target: job.CargoReference},
		UnitsCount:         int64(job.Weight / cargo.Mass),
		ShortestDistanceKM: job.Distance,

		// Some static data
		FillRatio: 1, // Dunno but other jobs have it at 1, so keep for now

		// Dunno where this data comes from, steal it from previous first job
		TrailerPlace: cTP,

		// Too lazy to implement, just steal it too
		CompanyTruck:      sii.Ptr{Target: cTruck},
		TrailerVariant:    sii.Ptr{Target: cTV},
		TrailerDefinition: sii.Ptr{Target: cTD},
	}
	j.Init("", jobID)

	// Add the new job to the game
	game.Entries = append(game.Entries, j)
	company.JobOffer = append([]sii.Ptr{{Target: j.Name()}}, company.JobOffer...)

	// Write the save-game
	if err = storeSave(vars["profileID"], storeSaveFolder, game, info); err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to store save"))
		return
	}

	apiGenericJSONResponse(w, http.StatusOK, map[string]interface{}{"success": true})
}

func handleFixPlayer(w http.ResponseWriter, r *http.Request) {
	var (
		fixType = fixAll
		vars    = mux.Vars(r)
	)

	if v := r.FormValue("type"); v != "" {
		fixType = v
	}

	game, info, err := loadSave(vars["profileID"], vars["saveFolder"])
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to load save"))
		return
	}

	info.SaveName = storeSaveName
	info.FileTime = time.Now().Unix()

	if err = fixPlayerTruck(game, fixType); err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to apply fixes"))
		return
	}

	if err = storeSave(vars["profileID"], storeSaveFolder, game, info); err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to store save"))
		return
	}

	apiGenericJSONResponse(w, http.StatusOK, map[string]interface{}{"success": true})
}

func handleGetSaveInfo(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)

	game, _, err := loadSave(vars["profileID"], vars["saveFolder"])
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to load save"))
		return
	}

	info, err := commSaveDetailsFromUnit(game)
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to gather info"))
		return
	}

	apiGenericJSONResponse(w, http.StatusOK, info)
}

func handleListJobs(w http.ResponseWriter, r *http.Request) {
	var (
		result       []commSaveJob
		vars         = mux.Vars(r)
		undiscovered = r.FormValue("undiscovered") == "true"
	)

	game, _, err := loadSave(vars["profileID"], vars["saveFolder"])
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to load save"))
		return
	}

	economy := game.BlocksByClass("economy")[0].(*sii.Economy)
	var visitedCities []string
	for _, p := range economy.VisitedCities {
		visitedCities = append(visitedCities, p.Target)
	}

	for _, cb := range game.BlocksByClass("company") {
		c := cb.(*sii.Company)

		cityName := strings.TrimPrefix(c.CityPtr().Target, "city.") // The "VisitedCities" pointers are kinda broken and do not contain the "city." part
		if !str.StringInSlice(cityName, visitedCities) && !undiscovered {
			continue
		}

		for _, b := range c.JobOffer {
			j := b.Resolve().(*sii.JobOfferData)

			if j.Target == "" || *j.ExpirationTime < economy.GameTime {
				continue
			}

			result = append(result, commSaveJob{
				OriginReference: c.Name(),
				TargetReference: strings.Join([]string{"company", "volatile", j.Target}, "."),
				CargoReference:  j.Cargo.Target,
				Distance:        j.ShortestDistanceKM,
				Urgency:         j.Urgency,
				Expires:         *j.ExpirationTime - economy.GameTime,
			})
		}
	}

	apiGenericJSONResponse(w, http.StatusOK, result)
}

func handleSetTrailer(w http.ResponseWriter, r *http.Request) {
	var (
		reference = r.FormValue("ref")
		vars      = mux.Vars(r)
	)

	game, info, err := loadSave(vars["profileID"], vars["saveFolder"])
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to load save"))
		return
	}

	if game.BlockByName(reference) == nil {
		apiGenericError(w, http.StatusBadRequest, errors.New("Invalid reference given"))
		return
	}

	info.SaveName = storeSaveName
	info.FileTime = time.Now().Unix()

	game.BlocksByClass("player")[0].(*sii.Player).AssignedTrailer = sii.Ptr{Target: reference}

	if err = storeSave(vars["profileID"], storeSaveFolder, game, info); err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to store save"))
		return
	}

	apiGenericJSONResponse(w, http.StatusOK, map[string]interface{}{"success": true})
}

func handleUpdateEconomyInfo(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)

	game, info, err := loadSave(vars["profileID"], vars["saveFolder"])
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to load save"))
		return
	}

	info.SaveName = storeSaveName
	info.FileTime = time.Now().Unix()

	blocks := game.BlocksByClass("economy")
	if len(blocks) != 1 {
		// expecting exactly one economy block
		apiGenericError(w, http.StatusInternalServerError, errors.New("Did not find economy block"))
		return
	}
	economy := blocks[0].(*sii.Economy)

	if xpRaw := r.FormValue("xp"); xpRaw != "" {
		xp, err := strconv.ParseInt(xpRaw, 10, 64)
		if err != nil {
			apiGenericError(w, http.StatusBadRequest, errors.Wrap(err, "Invalid value to xp parameter"))
			return
		}
		economy.ExperiencePoints = xp
	}

	if moneyRaw := r.FormValue("money"); moneyRaw != "" {
		money, err := strconv.ParseInt(moneyRaw, 10, 64)
		if err != nil {
			apiGenericError(w, http.StatusBadRequest, errors.Wrap(err, "Invalid value to money parameter"))
			return
		}
		economy.Bank.Resolve().(*sii.Bank).MoneyAccount = money
	}

	if err = storeSave(vars["profileID"], storeSaveFolder, game, info); err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to store save"))
		return
	}

	apiGenericJSONResponse(w, http.StatusOK, map[string]interface{}{"success": true})
}
