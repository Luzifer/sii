package main

import "net/http"

func init() {
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveName}", handleGetSaveInfo).Methods(http.MethodGet)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveName}/fix", handleFixPlayer).Methods(http.MethodPut)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveName}/jobs", handleListJobs).Methods(http.MethodGet)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveName}/jobs", handleAddJob).Methods(http.MethodPost)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveName}/jobs/{jobID}", handleEditJob).Methods(http.MethodPut)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveName}/set-trailer", handleSetTrailer).Methods(http.MethodPut)
}

func handleAddJob(w http.ResponseWriter, r *http.Request) {
	// FIXME: Implementation missing
}

func handleEditJob(w http.ResponseWriter, r *http.Request) {
	// FIXME: Implementation missing
}

func handleFixPlayer(w http.ResponseWriter, r *http.Request) {
	var fixType = "all"
	if v := r.FormValue("type"); v != "" {
		fixType = v
	}

	_ = fixType

	// FIXME: Implementation missing
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
