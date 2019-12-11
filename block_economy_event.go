package sii

func init() {
	RegisterBlock(&EconomyEvent{})
}

type EconomyEvent struct {
	Time     int64 `sii:"time"`
	UnitLink Ptr   `sii:"unit_link"`
	Param    int64 `sii:"param"`

	blockName string
}

func (EconomyEvent) Class() string { return "economy_event" }

func (e *EconomyEvent) Init(class, name string) {
	e.blockName = name
}

func (e EconomyEvent) Name() string { return e.blockName }
