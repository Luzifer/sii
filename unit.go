package sii

type Block interface {
	Name() string
	Type() string
}

type Marshaler interface {
	MarshalSII() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalSII([]byte) error
}

type RawBlock struct {
	Data []byte

	blockName string
	blockType string
}

func (r RawBlock) MarshalSII() []byte { return r.Data }
func (r RawBlock) Name() string       { return r.blockName }
func (r RawBlock) Type() string       { return r.blockType }
func (r *RawBlock) UnmarshalSII(blockName, blockType string, in []byte) error {
	r.Data = in
	return nil
}

type Unit struct {
	Entries []Block
}
