package sii

func init() {
	RegisterBlock(&EconomyEventQueue{})
}

type EconomyEventQueue struct {
	Data []Ptr `sii:"data"`

	blockName string
}

func (EconomyEventQueue) Class() string { return "economy_event_queue" }

func (e *EconomyEventQueue) Init(class, name string) {
	e.blockName = name
}

func (e EconomyEventQueue) Name() string { return e.blockName }
