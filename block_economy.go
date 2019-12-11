package sii

func init() {
	RegisterBlock(&Economy{})
}

type Economy struct {
	Bank                           Ptr        `sii:"bank"`
	Player                         Ptr        `sii:"player"`
	Companies                      []Ptr      `sii:"companies"`
	Garages                        []Ptr      `sii:"garages"`
	GarageIgnoreList               []Ptr      `sii:"garage_ignore_list"`
	GameProgress                   Ptr        `sii:"game_progress"`
	EventQueue                     Ptr        `sii:"event_queue"`
	MailCtrl                       Ptr        `sii:"mail_ctrl"`
	OversizeOfferCtrl              Ptr        `sii:"oversize_offer_ctrl"`
	GameTime                       int64      `sii:"game_time"`
	GameTimeSecs                   float32    `sii:"game_time_secs"`
	GameTimeInitial                int64      `sii:"game_time_initial"`
	AchievementsAdded              int64      `sii:"achievements_added"`
	NewGame                        bool       `sii:"new_game"`
	TotalDistance                  int64      `sii:"total_distance"`
	ExperiencePoints               int64      `sii:"experience_points"`
	ADR                            int64      `sii:"adr"` // Needs verification
	LongDist                       int64      `sii:"long_dist"`
	Heavy                          int64      `sii:"heavy"`
	Fragile                        int64      `sii:"fragile"`
	Urgent                         int64      `sii:"urgent"`
	Mechanical                     int64      `sii:"mechanical"`
	UserColors                     []int64    `sii:"user_colors"`
	DeliveryLog                    Ptr        `sii:"delivery_log"`
	FerryLog                       Ptr        `sii:"ferry_log"`
	StoredCameraMode               int64      `sii:"stored_camera_mode"`
	StoredActorState               int64      `sii:"stored_actor_state"`
	StoredHighBeamStyle            int64      `sii:"stored_high_beam_style"`
	StoredActorWiperMode           int64      `sii:"stored_actor_wiper_mode"`
	StoredActorRetarder            int64      `sii:"stored_actor_retarder"`
	StoredDisplayMode              int64      `sii:"stored_display_mode"`
	StoredDashboardMapMode         int64      `sii:"stored_dashboard_map_mode"`
	StoredWorldMapZoom             int64      `sii:"stored_world_map_zoom"`
	StoredOnlineJobID              int64      `sii:"stored_online_job_id"`
	StoredOnlineGPSBehind          int64      `sii:"stored_online_gps_behind"`
	StoredOnlineGPSAhead           int64      `sii:"stored_online_gps_ahead"`
	StoredOnlineGPSBehindWaypoints []Ptr      `sii:"stored_online_gps_behind_waypoints"`
	StoredOnlineGPSAheadWaypoints  []Ptr      `sii:"stored_online_gps_ahead_waypoints"`
	StoredOnlineGPSAvoidWaypoints  []Ptr      `sii:"stored_online_gps_avoid_waypoints"`
	StoredSpecialJob               Ptr        `sii:"stored_special_job"`
	PoliceCtrl                     Ptr        `sii:"police_ctrl"`
	StoredMapState                 int64      `sii:"stored_map_state"`
	StoredGasPumpMoney             int64      `sii:"stored_gas_pump_money"`
	StoredWeatherChangeTimer       float32    `sii:"stored_weather_change_timer"`
	StoredCurrentWeather           int64      `sii:"stored_current_weather"`
	StoredRainWetness              int64      `sii:"stored_rain_wetness"`
	TimeZone                       int64      `sii:"time_zone"`
	TimeZoneName                   string     `sii:"time_zone_name"`
	LastFerryPosition              [3]float32 `sii:"last_ferry_position"`
	StoredShowWeigh                bool       `sii:"stored_show_weigh"`
	StoredNeedToWeigh              bool       `sii:"stored_need_to_weigh"`
	StoredNavStartPos              [3]float32 `sii:"stored_nav_start_pos"`
	StoredNavEndPos                [3]float32 `sii:"stored_nav_end_pos"`
	StoredGPSBehind                int64      `sii:"stored_gps_behind"`
	StoredGPSAhead                 int64      `sii:"stored_gps_ahead"`
	StoredGPSBehindWaypoints       []Ptr      `sii:"stored_gps_behind_waypoints"`
	StoredGPSAheadWaypoints        []Ptr      `sii:"stored_gps_ahead_waypoints"`
	StoredGPSAvoidWaypoints        []Ptr      `sii:"stored_gps_avoid_waypoints"`
	StoredStartTollgatePos         [3]float32 `sii:"stored_start_tollgate_pos"`
	StoredTutorialState            int64      `sii:"stored_tutorial_state"`
	StoredMapActions               []Ptr      `sii:"stored_map_actions"`
	CleanDistanceCounter           int64      `sii:"clean_distance_counter"`
	CleanDistanceMax               int64      `sii:"clean_distance_max"`
	NoCargoDamageDistanceCounter   int64      `sii:"no_cargo_damage_distance_counter"`
	NoCargoDamageDistanceMax       int64      `sii:"no_cargo_damage_distance_max"`
	NoViolationDistanceCounter     int64      `sii:"no_violation_distance_counter"`
	NoViolationDistanceMax         int64      `sii:"no_violation_distance_max"`
	TotalRealTime                  int64      `sii:"total_real_time"`
	RealTimeSeconds                float32    `sii:"real_time_seconds"`
	VisitedCities                  []Ptr      `sii:"visited_cities"`
	VisitedCitiesCount             []int64    `sii:"visited_cities_count"`
	LastVisitedCity                Ptr        `sii:"last_visited_city"`
	UnlockedDealers                []Ptr      `sii:"unlocked_dealers"`
	UnlockedRecruitments           []Ptr      `sii:"unlocked_recruitments"`
	TotalScreenshotCount           int64      `sii:"total_screeshot_count"`
	UndamagedCargoRow              int64      `sii:"undamaged_cargo_row"`
	ServiceVisitCount              int64      `sii:"service_visit_count"`
	LastServicePos                 [3]float32 `sii:"last_service_pos"`
	GasStationVisitCount           int64      `sii:"gas_station_visit_count"`
	LastGasStationPos              [3]float32 `sii:"last_gas_station_pos"`
	EmergencyCallCount             int64      `sii:"emergency_call_count"`
	AICrashCount                   int64      `sii:"ai_crash_count"`
	TruckColorChangeCount          int64      `sii:"truck_color_change_count"`
	RedLightFineCount              int64      `sii:"red_light_fine_count"`
	CancelledJobCount              int64      `sii:"cancelled_job_count"`
	TotalFuelLitres                int64      `sii:"total_fuel_litres"`
	TotalFuelPrice                 int64      `sii:"total_fuel_price"`
	TransportedCargoTypes          []Ptr      `sii:"transported_cargo_types"`
	AchievedFeats                  []Ptr      `sii:"achieved_feats"` // Needs verification
	DiscoveredRoads                []Ptr      `sii:"discovered_roads"`
	DiscoveredItems                []int64    `sii:"discovered_items"` // Needs verification
	DriversOffer                   []Ptr      `sii:"drivers_offer"`
	FreelanceTruckOffer            Ptr        `sii:"freelance_truck_offer"`
	TrucksBoughtOnline             int64      `sii:"trucks_bought_online"`
	SpecialCargoTimer              int64      `sii:"special_cargo_timer"`
	ScreenAccessList               []string   `sii:"screen_access_list"`
	DriverPool                     []Ptr      `sii:"driver_pool"`
	Registry                       Ptr        `sii:"registry"`
	CompanyJobsInvitationSent      bool       `sii:"company_jobs_invitation_sent"`
	CompanyCheckHash               RawValue   `sii:"company_check_hash"` // Too long for int, not float, not string, wat?
	Relations                      []int64    `sii:"relations"`
	BusStops                       []Ptr      `sii:"bus_stops"`
	BusJobLog                      Ptr        `sii:"bus_job_log"`
	BusExperiencePoints            int64      `sii:"bus_experience_points"`
	BusTotalDistance               int64      `sii:"bus_total_distance"`
	BusFinishedJobCount            int64      `sii:"bus_finished_job_count"`
	BusCancelledJobCount           int64      `sii:"bus_cancelled_job_count"`
	BusTotalPassengers             int64      `sii:"bus_total_passengers"`
	BusTotalStops                  int64      `sii:"bus_total_stops"`
	BusGameTime                    int64      `sii:"bus_game_time"`
	BusPlayingTime                 int64      `sii:"bus_playing_time"`

	blockName string
}

func (Economy) Class() string { return "economy" }

func (e *Economy) Init(class, name string) {
	e.blockName = name
}

func (e Economy) Name() string { return e.blockName }
