package sii

type Block interface {
	Class() string
	Init(class, name string)
	Name() string
}

type Marshaler interface {
	MarshalSII() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalSII([]byte) error
}

type RawBlock struct {
	Data []byte

	blockName  string
	blockClass string
}

func (r *RawBlock) Init(class, name string) {
	r.blockClass = class
	r.blockName = name
}
func (r RawBlock) MarshalSII() ([]byte, error) { return r.Data, nil }
func (r RawBlock) Name() string                { return r.blockName }
func (r RawBlock) Class() string               { return r.blockClass }
func (r *RawBlock) UnmarshalSII(in []byte) error {
	r.Data = in
	return nil
}

type Unit struct {
	Entries []Block
}
