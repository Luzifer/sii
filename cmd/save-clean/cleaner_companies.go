package main

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/go_helpers/str"
	"github.com/Luzifer/sii"
)

func cleanCompanies(baseGameUnit, saveGame *sii.Unit) error {
	for _, b := range saveGame.BlocksByClass("company") {
		c, ok := b.(*sii.Company)
		if !ok {
			// Should not happen but to be sure...
			continue
		}

		var companyOK = true

		if _, ok := baseGameUnit.BlockByName(c.CityPtr().Target).(*sii.CityData); !ok {
			log.WithFields(log.Fields{
				"company": c.Name(),
				"reason":  "Missing city",
			}).Warn("Marking company invalid")
			companyOK = false
		}

		if _, ok := baseGameUnit.BlockByName(c.PermanentData.Target).(*sii.CompanyPermanent); !ok {
			log.WithFields(log.Fields{
				"company": c.Name(),
				"reason":  "Missing company permanent data",
			}).Warn("Marking company invalid")
			companyOK = false
		}

		if companyOK {
			continue
		}

		var relatedJobOfferData []string
		for _, p := range c.JobOffer {
			relatedJobOfferData = append(relatedJobOfferData, p.Target)
		}

		// Remove traces of broken company
		eco := saveGame.BlocksByClass("economy")[0].(*sii.Economy)

		// Remove blocks with references to this company (and the company itself
		var tmpBlocks []sii.Block
		for _, bl := range saveGame.Entries {
			switch bl.(type) {

			case *sii.Company:
				// Found the company, skip re-adding it
				if b == bl {
					continue
				}

			case *sii.EconomyEvent:
				// Event is a scheduled event for this company, do not re-add it
				if bl.(*sii.EconomyEvent).UnitLink.Target == c.Name() {
					continue
				}

			case *sii.JobOfferData:
				// If this offer starts at the company, skip re-adding it
				if str.StringInSlice(bl.Name(), relatedJobOfferData) {
					continue
				}

				// Job is directed to this company, do not re-add
				if strings.HasSuffix(c.Name(), bl.(*sii.JobOfferData).Target) {
					continue
				}

			}

			tmpBlocks = append(tmpBlocks, bl)
		}
		saveGame.Entries = tmpBlocks

		// Deregister company
		var tmpComps []sii.Ptr
		for _, cp := range eco.Companies {
			if cp.Target == c.Name() {
				// Skip this and therefore remove the reference
				continue
			}
			tmpComps = append(tmpComps, cp)
		}
		eco.Companies = tmpComps
	}

	return nil
}
