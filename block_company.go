package sii

import "strings"

func init() {
	RegisterBlock(&Company{})
}

type Company struct {
	PermanentData       Ptr         `sii:"permanent_data"` // external pointer
	DelieredTrailer     Ptr         `sii:"delivered_trailer"`
	DeliveredPos        []Placement `sii:"delivered_pos"`
	JobOffer            []Ptr       `sii:"job_offer"`
	CargoOfferSeeds     []int64     `sii:"cargo_offer_seeds"`
	Discovered          bool        `sii:"discovered"`
	ReservedTrailerSlot Ptr         `sii:"reserved_trailer_slot"` // TODO: Maybe wrong type, haven't seen other than "nil"

	blockName string
}

func (Company) Class() string { return "company" }

func (c *Company) Init(class, name string) {
	c.blockName = name
}

func (c Company) Name() string { return c.blockName }

func (c Company) CityPtr() *Ptr {
	nameParts := strings.Split(c.Name(), ".")
	if len(nameParts) != 4 || nameParts[0] != "company" || nameParts[1] != "volatile" {
		return nil
	}

	return &Ptr{Target: strings.Join([]string{"city", nameParts[3]}, ".")}
}
