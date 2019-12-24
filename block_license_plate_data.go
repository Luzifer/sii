package sii

func init() {
	RegisterBlock(&LicensePlateData{})
}

type LicensePlateData struct {
	Type Ptr `sii:"type"`

	Templates []string `sii:"templates"`

	Def0 []string `sii:"def0"`
	Def1 []string `sii:"def1"`
	Def2 []string `sii:"def2"`
	Def3 []string `sii:"def3"`

	blockName string
}

func (LicensePlateData) Class() string { return "license_plate_data" }

func (l *LicensePlateData) Init(class, name string) {
	l.blockName = name
}

func (l LicensePlateData) Name() string { return l.blockName }
