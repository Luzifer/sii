package sii

func init() {
	RegisterBlock(&CityData{})
}

type CityData struct {
	CityName          string `sii:"city_name"`
	CityNameLocalized string `sii:"city_name_localized"`
	Country           Ptr    `sii:"country"`

	MapXOffsets []int64 `sii:"map_x_offsets"`
	MapYOffsets []int64 `sii:"map_y_offsets"`

	VehicleBrands []string `sii:"vehicle_brands"`

	LicensePlate []Ptr `sii:"license_plate"`

	blockName string
}

func (CityData) Class() string { return "city_data" }

func (c *CityData) Init(class, name string) {
	c.blockName = name
}

func (c CityData) Name() string { return c.blockName }
