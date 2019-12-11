package sii

func init() {
	RegisterBlock(&Registry{})
}

type Registry struct {
	Data  []int64 `sii:"data"`
	Valid []bool  `sii:"valid"`
	Keys  []int64 `sii:"keys"`
	Index []int64 `sii:"index"`

	blockName string
}

func (Registry) Class() string { return "registry" }

func (r *Registry) Init(class, name string) {
	r.blockName = name
}

func (r Registry) Name() string { return r.blockName }
