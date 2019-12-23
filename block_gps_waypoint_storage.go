package sii

func init() {
	RegisterBlock(&GPSWaypointStorage{})
}

type GPSWaypointStorage struct {
	NavNodePosition [3]int64 `sii:"nav_node_position"`
	Direction       Ptr      `sii:"direction"`

	blockName string
}

func (GPSWaypointStorage) Class() string { return "gps_waypoint_storage" }

func (g *GPSWaypointStorage) Init(class, name string) {
	g.blockName = name
}

func (g GPSWaypointStorage) Name() string { return g.blockName }
