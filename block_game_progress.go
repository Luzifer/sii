package sii

func init() {
	RegisterBlock(&GameProgress{})
}

type GameProgress struct {
	GenericTransports   Ptr      `sii:"generic_transports"`
	UndamagedTransports Ptr      `sii:"undamaged_transports"`
	CleanTransports     Ptr      `sii:"clean_transports"`
	OwnedTrucks         []string `sii:"owned_trucks"`

	blockName string
}

func (GameProgress) Class() string { return "game_progress" }

func (g *GameProgress) Init(class, name string) {
	g.blockName = name
}

func (g GameProgress) Name() string { return g.blockName }
