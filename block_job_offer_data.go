package sii

func init() {
	RegisterBlock(&JobOfferData{})
}

type JobOfferData struct {
	Target             string  `sii:"target"`
	ExpirationTime     int64   `sii:"expiration_time"`
	Urgency            int64   `sii:"urgency"`
	ShortestDistanceKM int64   `sii:"shortest_distance_km"`
	FerryTime          int64   `sii:"ferry_time"`
	FerryPrice         int64   `sii:"ferry_price"`
	Cargo              Ptr     `sii:"cargo"`              // External pointer
	CompanyTruck       Ptr     `sii:"company_truck"`      // Partial external pointer?
	TrailerVariant     Ptr     `sii:"trailer_variant"`    // External pointer
	TrailerDefinition  Ptr     `sii:"trailer_definition"` // External pointer
	UnitsCount         int64   `sii:"units_count"`
	FillRatio          float32 `sii:"fill_ratio"`
	TrailerPlace       int64   `sii:"trailer_place"` // TODO: What's this? Seems to be "0" in all jobs

	blockName string
}

func (JobOfferData) Class() string { return "job_offer_data" }

func (j *JobOfferData) Init(class, name string) {
	j.blockName = name
}

func (j JobOfferData) Name() string { return j.blockName }
