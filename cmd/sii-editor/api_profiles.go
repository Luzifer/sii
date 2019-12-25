package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const savePollTime = 10 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

	// Remember latest save
	var latestSave time.Time
	for _, s := range saves {
		if s.SaveTime.After(latestSave) {
			latestSave = s.SaveTime
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Debug("Unable to open websocket")
		return
	}
	defer conn.Close()

	if err = conn.WriteJSON(saves); err != nil {
		log.WithError(err).Error("Unable to send saves list")
		return
	}

	for t := time.NewTicker(savePollTime); ; <-t.C {
		saves, err = listSaves(vars["profileID"])
		if err != nil {
			log.WithError(err).Error("Unable to list saves during socket")
			return
		}

		var newSaveTime time.Time
		for _, s := range saves {
			if s.SaveTime.After(latestSave) {
				newSaveTime = s.SaveTime
			}
		}

		if newSaveTime.IsZero() {
			// Nothing new
			continue
		}

		if err = conn.WriteJSON(saves); err != nil {
			log.WithError(err).Error("Unable to send saves list")
			return
		}

		latestSave = newSaveTime
	}
}

func handleListProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := listProfiles()
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, err)
		return
	}

	apiGenericJSONResponse(w, http.StatusOK, profiles)
}
