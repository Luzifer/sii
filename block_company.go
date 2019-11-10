package sii

func init() {
	RegisterBlock(&Company{})
}

type Company struct {
	PermanentData       Ptr     `sii:"permanent_data"` // external pointer
	DelieredTrailer     Ptr     `sii:"delivered_trailer"`
	DeliveredPos        int64   `sii:"delivered_pos"`
	JobOffer            []Ptr   `sii:"job_offer"`
	CargoOfferSeeds     []int64 `sii:"cargo_offer_seeds"`
	Discovered          bool    `sii:"discovered"`
	ReservedTrailerSlot int64   `sii:"reserved_trailer_slot"` // TODO: Maybe wrong type, haven't seen other than "nil"

	blockName string
}

func (Company) Class() string { return "company" }

func (c *Company) Init(class, name string) {
	c.blockName = name
}

func (c Company) Name() string { return c.blockName }
