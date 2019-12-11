package sii

func init() {
	RegisterBlock(&FerryLog{})
}

type FerryLog struct {
	Entries []Ptr `sii:"entries"`

	blockName string
}

func (FerryLog) Class() string { return "ferry_log" }

func (f *FerryLog) Init(class, name string) {
	f.blockName = name
}

func (f FerryLog) Name() string { return f.blockName }
