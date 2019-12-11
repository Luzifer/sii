package sii

func init() {
	RegisterBlock(&JobInfo{})
}

type JobInfo struct {
	Cargo             Ptr     `sii:"cargo"`
	SourceCompany     Ptr     `sii:"source_company"`
	TargetCompany     Ptr     `sii:"target_company"`
	IsArticulated     bool    `sii:"is_articulated"`
	IsCargoMarketJob  bool    `sii:"is_cargo_market_job"`
	StartTime         int64   `sii:"start_time"`
	PlannedDistanceKM int64   `sii:"planned_distance_km"`
	FerryTime         int64   `sii:"ferry_time"`
	FerryPrice        int64   `sii:"ferry_price"`
	Urgency           *int64  `sii:"urgency"`
	Special           Ptr     `sii:"special"`
	UnitCount         int64   `sii:"units_count"`
	FillRatio         float32 `sii:"fill_ratio"`

	blockName string
}

func (JobInfo) Class() string { return "job_info" }

func (j *JobInfo) Init(class, name string) {
	j.blockName = name
}

func (j JobInfo) Name() string { return j.blockName }
