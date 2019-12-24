package main

import (
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
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/jobs/{jobID}", handleEditJob).Methods(http.MethodPut)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/set-trailer", handleSetTrailer).Methods(http.MethodPut)
}

func handleAddJob(w http.ResponseWriter, r *http.Request) {
	// FIXME: Implementation missing
}

func handleEditJob(w http.ResponseWriter, r *http.Request) {
	// FIXME: Implementation missing
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

			if j.Target == "" {
				continue
			}

			result = append(result, commSaveJob{
				OriginReference: c.Name(),
				TargetReference: strings.Join([]string{"company", "volatile", j.Target}, "."),
				CargoReference:  j.Cargo.Target,
				Distance:        j.ShortestDistanceKM,
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
