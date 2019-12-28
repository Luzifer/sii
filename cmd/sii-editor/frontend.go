package main

import "net/http"

func init() {
	router.HandleFunc("/", handleIndexPage).Methods(http.MethodGet)
	router.PathPrefix("/asset/").Handler(
		http.StripPrefix("/asset/", http.FileServer(http.Dir("frontend"))),
	).Methods(http.MethodGet)
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/index.html")
}
