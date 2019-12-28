package main

import (
	"encoding/json"
	"net/http"
)

func apiGenericError(w http.ResponseWriter, status int, err error) {
	var eString = "undefined error"
	if err != nil {
		eString = err.Error()
	}

	data := map[string]interface{}{
		"code":    status,
		"error":   eString,
		"success": false,
	}

	apiGenericJSONResponse(w, status, data)
}

func apiGenericJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}
