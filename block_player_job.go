package sii

func init() {
	RegisterBlock(&PlayerJob{})
}

type PlayerJob struct {
	CompanyTruck          Ptr       `sii:"company_truck"`
	CompanyTrailer        Ptr       `sii:"company_trailer"`
	TargetPlacement       Placement `sii:"target_placement"`
	TargetPlacementMedium Placement `sii:"target_placement_medium"`
	TargetPlacementHard   Placement `sii:"target_placement_hard"`
	TargetPlacementRigid  Placement `sii:"target_placement_rigid"`
	SourcePlacement       Placement `sii:"source_placement"`
	SelectedTarget        int64     `sii:"selected_target"` // Needs verification
	TimeLowerLimit        int64     `sii:"time_lower_limit"`
	TimeUpperLimit        int64     `sii:"time_upper_limit"`
	JobDistance           int64     `sii:"job_distance"`
	FuelConsumed          float32   `sii:"fuel_consumed"`
	LastReportedFuel      float32   `sii:"last_reported_fuel"`
	TotalFines            int64     `sii:"total_fines"` // Needs verification
	IsTrailerLoaded       bool      `sii:"is_trailer_loaded"`
	OnlineJobID           *int64    `sii:"online_job_id"` // Needs verification
	OnlineJobTrailerModel Ptr       `sii:"online_job_trailer_model"`
	AutoloadUsed          bool      `sii:"autoload_used"`
	Cargo                 Ptr       `sii:"cargo"` // External pointer
	SourceCompany         Ptr       `sii:"source_company"`
	TargetCompany         Ptr       `sii:"target_company"`
	IsArticulated         bool      `sii:"is_articulated"`
	IsCargoMarketJob      bool      `sii:"is_cargo_market_job"`
	StartTime             int64     `sii:"start_time"`
	PlannedDistanceKM     int64     `sii:"planned_distance_km"`
	FerryTime             int64     `sii:"ferry_time"`
	FerryPrice            int64     `sii:"ferry_price"`
	Urgency               *int64    `sii:"urgency"`
	Special               Ptr       `sii:"special"`
	UnitCount             int64     `sii:"units_count"`
	FillRatio             float32   `sii:"fill_ratio"`

	blockName string
}

func (PlayerJob) Class() string { return "player_job" }

func (p *PlayerJob) Init(class, name string) {
	p.blockName = name
}

func (p PlayerJob) Name() string { return p.blockName }
