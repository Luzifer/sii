package sii

func init() {
	RegisterBlock(&Vehicle{})
}

type Vehicle struct {
	FuelRelative              float32      `sii:"fuel_relative"`
	RheostatFactor            float32      `sii:"rheostat_factor"`
	UserMirrorRot             [][4]float32 `sii:"user_mirror_rot"`
	UserHeadOffset            [3]float32   `sii:"user_head_offset"`
	UserFOV                   int64        `sii:"user_fov"`
	UserWheelUpDown           int64        `sii:"user_wheel_up_down"`
	UserWheelFrontBack        int64        `sii:"user_wheel_front_back"`
	UserMouseLeftRightDefault int64        `sii:"user_mouse_left_right_default"`
	UserMouseUpDownDefault    int64        `sii:"user_mouse_up_down_default"`
	Accessories               []Ptr        `sii:"accessories"`
	Odometer                  int64        `sii:"odometer"`
	OdometerFloatPart         float32      `sii:"odometer_float_part"`
	TripFuelL                 int64        `sii:"trip_fuel_l"`
	TripFuel                  float32      `sii:"trip_fuel"`
	TripDistanceKM            int64        `sii:"trip_distance_km"`
	TripDistance              float32      `sii:"trip_distance"`
	LicensePlate              string       `sii:"license_plate"`

	blockName string
}

func (Vehicle) Class() string { return "vehicle" }

func (v *Vehicle) Init(class, name string) {
	v.blockName = name
}

func (v Vehicle) Name() string { return v.blockName }
