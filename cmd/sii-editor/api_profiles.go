package main

import (
	"encoding/json"
	"net/http"
)

func init() {
	router.HandleFunc("/api/profiles", handleListProfiles).Methods(http.MethodGet)
	router.HandleFunc("/api/profiles/{profileID}/saves", handleGetProfileSaves).Methods(http.MethodGet)
}

func handleGetProfileSaves(w http.ResponseWriter, r *http.Request) {
	var subscribe = r.FormValue("subscribe") == "true"

	_ = subscribe // If so open socket and let browser know there are new saves

	// FIXME: Implementation missing
}

func handleListProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := listProfiles()
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, err)
		return
	}

	json.NewEncoder(w).Encode(profiles)
}
