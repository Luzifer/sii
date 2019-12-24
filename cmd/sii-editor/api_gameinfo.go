package main

import (
	"net/http"

	"github.com/Luzifer/sii"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func init() {
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/companies", handleListCompanies).Methods(http.MethodGet)
}

type commCompany struct {
	City string
	Name string
}

func handleListCompanies(w http.ResponseWriter, r *http.Request) {
	var (
		result = map[string]commCompany{}
		vars   = mux.Vars(r)
	)

	game, _, err := loadSave(vars["profileID"], vars["saveFolder"])
	if err != nil {
		apiGenericError(w, http.StatusInternalServerError, errors.Wrap(err, "Unable to load save"))
		return
	}

	for _, b := range game.BlocksByClass("company") {
		c, ok := b.(*sii.Company)
		if !ok {
			// Should not happen but to be sure...
			continue
		}

		result[c.Name()] = commCompany{
			City: baseGameUnit.BlockByName(c.CityPtr().Target).(*sii.CityData).CityName,
			Name: baseGameUnit.BlockByName(c.PermanentData.Target).(*sii.CompanyPermanent).CompanyName,
		}
	}

	apiGenericJSONResponse(w, http.StatusOK, result)
}
