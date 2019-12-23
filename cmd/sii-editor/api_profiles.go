package main

import "net/http"

func init() {
	router.HandleFunc("/api/profiles", handleListProfiles).Methods(http.MethodGet)
	router.HandleFunc("/api/profiles/{profileID}", handleGetProfileInfo).Methods(http.MethodGet)
}

func handleGetProfileInfo(w http.ResponseWriter, r *http.Request) {
	// FIXME: Implementation missing
}

func handleListProfiles(w http.ResponseWriter, r *http.Request) {
	// FIXME: Implementation missing
}
