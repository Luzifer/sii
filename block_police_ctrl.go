package sii

func init() {
	RegisterBlock(&PoliceCtrl{})
}

type PoliceCtrl struct {
	OffenceTimer   []float32 `sii:"offence_timer"`
	OffenceCounter []int64   `sii:"offence_counter"`
	OffenceValid   []bool    `sii:"offence_valid"`

	blockName string
}

func (PoliceCtrl) Class() string { return "police_ctrl" }

func (p *PoliceCtrl) Init(class, name string) {
	p.blockName = name
}

func (p PoliceCtrl) Name() string { return p.blockName }
