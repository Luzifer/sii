package sii

func init() {
	RegisterBlock(&TrailerUtilizationLog{})
}

type TrailerUtilizationLog struct {
	Entries                 []Ptr   `sii:"entries"` // Needs verification
	TotalDrivenDistanceKM   int64   `sii:"total_driven_distance_km"`
	TotalTransportedCargoes int64   `sii:"total_transported_cargoes"`
	TotalTransportedWeight  float32 `sii:"total_transported_weight"`

	blockName string
}

func (TrailerUtilizationLog) Class() string { return "trailer_utilization_log" }

func (t *TrailerUtilizationLog) Init(class, name string) {
	t.blockName = name
}

func (t TrailerUtilizationLog) Name() string { return t.blockName }
