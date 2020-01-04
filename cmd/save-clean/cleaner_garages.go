package main

import (
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/go_helpers/str"
	"github.com/Luzifer/sii"
)

func cleanGarages(baseGameUnit, saveGame *sii.Unit) error {
	var (
		deletions                             []string
		eco                                   = saveGame.BlocksByClass("economy")[0].(*sii.Economy)
		freeDrivers, freeTrailers, freeTrucks []sii.Ptr
	)

	for _, b := range saveGame.BlocksByClass("garage") {
		g, ok := b.(*sii.Garage)
		if !ok {
			// Should not happen but to be sure...
			continue
		}

		var cityRef = strings.Join([]string{"city", strings.TrimPrefix(g.Name(), "garage.")}, ".")
		if _, ok := baseGameUnit.BlockByName(cityRef).(*sii.CityData); ok {
			continue
		}

		log.WithFields(log.Fields{
			"garage": g.Name(),
			"reason": "Missing city",
		}).Warn("Marking garage invalid")

		// Don't trash trucks / trailers / drivers
		for _, p := range g.Vehicles {
			if !p.IsNull() {
				freeTrucks = append(freeTrucks, p)
			}
		}

		for _, p := range g.Trailers {
			if !p.IsNull() {
				freeTrailers = append(freeTrailers, p)
			}
		}

		for _, p := range g.Drivers {
			if !p.IsNull() {
				freeDrivers = append(freeDrivers, p)
			}
		}

		var tmpGarages []sii.Ptr
		for _, p := range eco.Garages {
			if p.Target == g.Name() {
				continue
			}
			tmpGarages = append(tmpGarages, p)
		}
		eco.Garages = tmpGarages

		deletions = append(deletions, g.Name())

		profitLog := g.ProfitLog.Resolve().(*sii.ProfitLog)
		deletions = append(deletions, profitLog.Name())

		for _, p := range profitLog.StatsData {
			if p.IsNull() {
				continue
			}
			deletions = append(deletions, p.Target)
		}

	}

	log.WithFields(log.Fields{
		"drivers":  len(freeDrivers),
		"trailers": len(freeTrailers),
		"trucks":   len(freeTrucks),
	}).Info("Reassignments needed")

	if c := len(freeDrivers) + len(freeTrailers) + len(freeTrucks); c > 0 {
		return errors.Errorf("%d reassignments needed, not yet supported")
	}

	var tmpBlocks []sii.Block
	for _, b := range saveGame.Entries {
		if str.StringInSlice(b.Name(), deletions) {
			// On list for deletion
			continue
		}
		tmpBlocks = append(tmpBlocks, b)
	}
	saveGame.Entries = tmpBlocks

	return nil
}
