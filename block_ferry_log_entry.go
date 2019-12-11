package sii

func init() {
	RegisterBlock(&FerryLogEntry{})
}

type FerryLogEntry struct {
	Ferry      Ptr   `sii:"ferry"`
	Connection Ptr   `sii:"connection"`
	LastVisit  int64 `sii:"last_visit"`
	UseCount   int64 `sii:"use_count"`

	blockName string
}

func (FerryLogEntry) Class() string { return "ferry_log_entry" }

func (f *FerryLogEntry) Init(class, name string) {
	f.blockName = name
}

func (f FerryLogEntry) Name() string { return f.blockName }
