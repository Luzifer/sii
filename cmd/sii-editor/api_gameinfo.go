package main

import (
	"net/http"
	"strings"

	"github.com/Luzifer/sii"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func init() {
	router.HandleFunc("/api/gameinfo/cargo", handleListCargo).Methods(http.MethodGet)
	router.HandleFunc("/api/profiles/{profileID}/saves/{saveFolder}/companies", handleListCompanies).Methods(http.MethodGet)
}

type commCargo struct {
	Name string  `json:"name"`
	Mass float32 `json:"mass"`
}

type commCompany struct {
	City string `json:"city"`
	Name string `json:"name"`
}

func handleListCargo(w http.ResponseWriter, r *http.Request) {
	var result = map[string]commCargo{}

	for _, b := range baseGameUnit.BlocksByClass("cargo_data") {
		c := b.(*sii.CargoData)
		var cName = c.CargoName
		if strings.HasPrefix(cName, "@@") {
			// Localization string, translate
			cName = locale.GetTranslation(cName)
		}

		result[c.Name()] = commCargo{
			Name: cName,
			Mass: c.Mass,
		}
	}

	apiGenericJSONResponse(w, http.StatusOK, result)
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
