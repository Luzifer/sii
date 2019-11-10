package sii

func init() {
	RegisterBlock(&DeliveryLog{})
	RegisterBlock(&DeliveryLogEntry{})
}

type DeliveryLog struct {
	Version        int   `sii:"version"`
	Entries        []Ptr `sii:"entries"`
	CachedJobCount int64 `sii:"cached_jobs_count"`

	blockName string
}

func (DeliveryLog) Class() string { return "delivery_log" }

func (d *DeliveryLog) Init(class, name string) {
	d.blockName = name
}

func (d DeliveryLog) Name() string { return d.blockName }

type DeliveryLogEntry struct {
	Params []RawValue `sii:"params"`

	blockName string
}

func (DeliveryLogEntry) Class() string { return "delivery_log_entry" }

func (d *DeliveryLogEntry) Init(class, name string) {
	d.blockName = name
}

func (d DeliveryLogEntry) Name() string { return d.blockName }
