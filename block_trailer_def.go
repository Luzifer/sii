package sii

func init() {
	RegisterBlock(&TrailerDef{})
}

type TrailerDef struct {
	Trailer                 string    `sii:"trailer"`
	GrossTrailerWeightLimit int64     `sii:"gross_trailer_weight_limit"`
	ChassisMass             int64     `sii:"chassis_mass"`
	BodyMass                int64     `sii:"body_mass"`
	Axles                   int64     `sii:"axles"`
	Volume                  int64     `sii:"volume"`
	BodyType                Ptr       `sii:"body_type"`
	ChainType               Ptr       `sii:"chain_type"`
	CountryValidity         int64     `sii:"country_validity"` // Needs verification
	MassRatio               []float32 `sii:"mass_ratio"`
	Length                  float32   `sii:"length"`
	SourceName              string    `sii:"source_name"`

	blockName string
}

func (TrailerDef) Class() string { return "trailer_def" }

func (t *TrailerDef) Init(class, name string) {
	t.blockName = name
}

func (t TrailerDef) Name() string { return t.blockName }
