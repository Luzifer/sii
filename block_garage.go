package sii

func init() {
	RegisterBlock(&Garage{})
}

type Garage struct {
	Vehicles     []Ptr   `sii:"vehicles"`
	Drivers      []Ptr   `sii:"drivers"`
	Trailers     []Ptr   `sii:"trailers"`
	Status       int     `sii:"status"`
	ProfitLog    Ptr     `sii:"profit_log"`
	Productivity float32 `sii:"productivity"`

	blockName string
}

func (Garage) Class() string { return "garage" }

func (g *Garage) Init(class, name string) {
	g.blockName = name
}

func (g Garage) Name() string { return g.blockName }
