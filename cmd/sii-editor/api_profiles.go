package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	router.HandleFunc("/api/profiles", handleListProfiles).Methods(http.MethodGet)
	router.HandleFunc("/api/profiles/{profileID}/saves", handleGetProfileSaves).Methods(http.MethodGet)
}

func handleGetProfileSaves(w http.ResponseWriter, r *http.Request) {
	var (
		subscribe = r.FormValue("subscribe") == "true"
		vars      = mux.Vars(r)
	)

	saves, err := listSaves(vars["profileID"])
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, err)
		return
	}

	if !subscribe {
		apiGenericJSONResponse(w, http.StatusOK, saves)
		return
	}

	_ = subscribe // If so open socket and let browser know there are new saves

	// FIXME: Implementation missing
}

func handleListProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := listProfiles()
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, err)
		return
	}

	apiGenericJSONResponse(w, http.StatusOK, profiles)
}
