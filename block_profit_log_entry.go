package sii

func init() {
	RegisterBlock(&ProfitLogEntry{})
}

type ProfitLogEntry struct {
	Revenue            int64 `sii:"revenue"`
	Wage               int64 `sii:"wage"`
	Maintenance        int64 `sii:"maintenance"`
	Fuel               int64 `sii:"fuel"`
	Distance           int64 `sii:"distance"`
	DistanceOnJob      bool  `sii:"distance_on_job"`
	CargoCount         int64 `sii:"cargo_count"`
	Cargo              Ptr   `sii:"cargo"`
	SourceCity         Ptr   `sii:"source_city"`
	SourceCompany      Ptr   `sii:"source_company"`
	DestinationCity    Ptr   `sii:"destination_city"`
	DestinationCompany Ptr   `sii:"destination_company"`
	TimestampDay       int64 `sii:"timestamp_day"`

	blockName string
}

func (ProfitLogEntry) Class() string { return "profit_log_entry" }

func (p *ProfitLogEntry) Init(class, name string) {
	p.blockName = name
}

func (p ProfitLogEntry) Name() string { return p.blockName }
