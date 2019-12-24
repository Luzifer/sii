package sii

func init() {
	RegisterBlock(&CargoData{})
}

type CargoData struct {
	CargoName       string  `sii:"name"`
	Fragility       float32 `sii:"fragility"`
	Groups          []Ptr   `sii:"group"`
	Volume          float32 `sii:"volume"`
	Mass            float32 `sii:"mass"`
	UnitRewardPerKM float32 `sii:"unit_reward_per_km"`
	UnitLoadTime    int64   `sii:"unit_load_time"`
	BodyTypes       []Ptr   `sii:"body_types"`

	blockName string
}

func (CargoData) Class() string { return "cargo_data" }

func (c *CargoData) Init(class, name string) {
	c.blockName = name
}

func (c CargoData) Name() string { return c.blockName }
