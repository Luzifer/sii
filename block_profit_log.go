package sii

func init() {
	RegisterBlock(&ProfitLog{})
}

type ProfitLog struct {
	StatsData        []Ptr  `sii:"stats_data"`
	AccDistanceFree  int64  `sii:"acc_distance_free"`
	AccDistanceOnJob int64  `sii:"acc_distance_on_job"`
	HistoryAge       *int64 `sii:"history_age"`

	blockName string
}

func (ProfitLog) Class() string { return "profit_log" }

func (p *ProfitLog) Init(class, name string) {
	p.blockName = name
}

func (p ProfitLog) Name() string { return p.blockName }
