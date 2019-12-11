package sii

func init() {
	RegisterBlock(&TransportData{})
}

type TransportData struct {
	Distance     int64    `sii:"distance"`
	Time         int64    `sii:"time"`
	Money        int64    `sii:"money"`
	CountPerADR  []int64  `sii:"count_per_adr"`
	Docks        []string `sii:"docks"`
	CountPerDock []int64  `sii:"count_per_dock"`

	blockName string
}

func (TransportData) Class() string { return "transport_data" }

func (t *TransportData) Init(class, name string) {
	t.blockName = name
}

func (t TransportData) Name() string { return t.blockName }
