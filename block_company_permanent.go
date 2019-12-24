package sii

func init() {
	RegisterBlock(&CompanyPermanent{})
}

type CompanyPermanent struct {
	CompanyName string `sii:"name"`
	SortName    string `sii:"sort_name"`
	TrailerLook Ptr    `sii:"trailer_look"`

	blockName string
}

func (CompanyPermanent) Class() string { return "company_permanent" }

func (c *CompanyPermanent) Init(class, name string) {
	c.blockName = name
}

func (c CompanyPermanent) Name() string { return c.blockName }
