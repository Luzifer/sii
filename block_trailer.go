package sii

import (
	"regexp"
	"strings"
)

func init() {
	RegisterBlock(&Trailer{})
}

type Trailer struct {
	TrailerDefinition       Ptr     `sii:"trailer_definition"`
	CargoMass               float32 `sii:"cargo_mass"`
	CargoDamage             float32 `sii:"cargo_damage"`               // Needs verification
	VirtualRearWheelsOffset int64   `sii:"virtual_rear_wheels_offset"` // Needs verification
	SlaveTrailer            Ptr     `sii:"slave_trailer"`
	IsPrivate               bool    `sii:"is_private"`
	Accessories             []Ptr   `sii:"accessories"`
	Odometer                int64   `sii:"odometer"`
	OdometerFloatPart       float32 `sii:"odometer_float_part"`
	TripFuelL               int64   `sii:"trip_fuel_l"`      // Needs verification
	TripFuel                int64   `sii:"trip_fuel"`        // Needs verification
	TripDistanceKM          int64   `sii:"trip_distance_km"` // Needs verification
	TripDistance            int64   `sii:"trip_distance"`    // Needs verification
	LicensePlate            string  `sii:"license_plate"`

	blockName string
}

func (Trailer) Class() string { return "trailer" }

func (t *Trailer) Init(class, name string) {
	t.blockName = name
}

func (t Trailer) Name() string { return t.blockName }

func (t Trailer) CleanedLicensePlate() string {
	return regexp.MustCompile(` +`).ReplaceAllString(
		regexp.MustCompile(`<[^>]+>`).ReplaceAllString(
			strings.Split(t.LicensePlate, "|")[0],
			" ",
		),
		" ",
	)
}
