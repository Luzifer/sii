package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const (
	storeSaveFolder = "sii_editor"
	storeSaveName   = "SII Editor"
)

func init() {
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}", handleGetSaveInfo).Methods(http.MethodGet)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/economy/{type}", handleUpdateEconomyInfo).Methods(http.MethodPut)
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
	// FIXME: Implementation missing
}

func handleListJobs(w http.ResponseWriter, r *http.Request) {
	// FIXME: Implementation missing
}

func handleSetTrailer(w http.ResponseWriter, r *http.Request) {
	// FIXME: Implementation missing
}

func handleUpdateEconomyInfo(w http.ResponseWriter, r *http.Request) {
	// FIXME: Implementation missing
}
