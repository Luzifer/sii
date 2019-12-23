package sii

func init() {
	RegisterBlock(&VehicleWheelAccessory{})
}

type VehicleWheelAccessory struct {
	Offset     int64      `sii:"offset"`
	PaintColor [3]float32 `sii:"paint_color"`
	Wear       float32    `sii:"wear"`
	DataPath   string     `sii:"data_path"`

	blockName string
}

func (VehicleWheelAccessory) Class() string { return "vehicle_wheel_accessory" }

func (v *VehicleWheelAccessory) Init(class, name string) {
	v.blockName = name
}

func (v VehicleWheelAccessory) Name() string { return v.blockName }
