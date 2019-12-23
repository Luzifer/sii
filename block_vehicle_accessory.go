package sii

func init() {
	RegisterBlock(&VehicleAccessory{})
}

type VehicleAccessory struct {
	Wear     float32 `sii:"wear"`
	DataPath string  `sii:"data_path"`

	blockName string
}

func (VehicleAccessory) Class() string { return "vehicle_accessory" }

func (v *VehicleAccessory) Init(class, name string) {
	v.blockName = name
}

func (v VehicleAccessory) Name() string { return v.blockName }
