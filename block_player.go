package sii

func init() {
	RegisterBlock(&Player{})
}

type Player struct {
	HQCity                   string      `sii:"hq_city"`
	Trailers                 []Ptr       `sii:"trailers"`
	TrailerUtilizationLogs   []Ptr       `sii:"trailer_utilization_logs"`
	TrailerDefs              []Ptr       `sii:"trailer_defs"`
	AssignedTruck            Ptr         `sii:"assigned_truck"`
	MyTruck                  Ptr         `sii:"my_truck"`
	MyTruckPlacement         Placement   `sii:"my_truck_placement"`
	MyTruckPlacementValid    bool        `sii:"my_truck_placement_valid"`
	MyTrailerPlacement       Placement   `sii:"my_trailer_placement"`
	MySlaveTrailerPlacements []Placement `sii:"my_slave_trailer_placements"`
	MyTrailerAttached        bool        `sii:"my_trailer_attached"`
	MyTrailerUsed            bool        `sii:"my_trailer_used"`
	AssignedTrailer          Ptr         `sii:"assigned_trailer"`
	MyTrailer                Ptr         `sii:"my_trailer"`
	AssignedTrailerConnected bool        `sii:"assigned_trailer_connected"`
	TruckPlacement           Placement   `sii:"truck_placement"`
	TrailerPlacement         Placement   `sii:"trailer_placement"`
	SlaveTrailerPlacements   []Placement `sii:"slave_trailer_placements"`
	ScheduleTransferToHQ     bool        `sii:"schedule_transfer_to_hq"`
	Flags                    uint64      `sii:"flags"` // ????
	GasPumpMoneyDebt         int64       `sii:"gas_pump_money_debt"`
	CurrentJob               Ptr         `sii:"current_job"`
	CurrentBusJob            Ptr         `sii:"current_bus_job"`
	SelectedJob              Ptr         `sii:"selected_job"`
	DrivingTime              int64       `sii:"driving_time"`
	SleepingCount            int         `sii:"sleeping_count"`
	FreeRoamDistance         int64       `sii:"free_roam_distance"`
	DiscoveryDistance        float32     `sii:"discovary_distance"` // Typo is intended and copied from real save-game
	DismissedDrivers         int         `sii:"dismissed_drivers"`
	Trucks                   []Ptr       `sii:"trucks"`
	TruckProfitLogs          []Ptr       `sii:"truck_profit_logs"`
	Drivers                  []Ptr       `sii:"drivers"`
	DriverReadinessTimer     []int64     `sii:"driver_readiness_timer"`
	DriverQuitWarned         []bool      `sii:"driver_quit_warned"`

	blockName string
}

func (Player) Class() string { return "player" }

func (p *Player) Init(class, name string) {
	p.blockName = name
}

func (p Player) Name() string { return p.blockName }
