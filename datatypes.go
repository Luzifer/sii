package sii

import (
	"bytes"
	"regexp"

	"github.com/pkg/errors"
)

// See https://modding.scssoft.com/wiki/Documentation/Engine/Units

// string => native type string

// float => native type float

// float2-4 => [2]float - [4]float

var placementRegexp = regexp.MustCompile(`^\(([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+)\) \(([0-9.-]+|&[0-9a-f]+); ([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+)\)$`)

// Placement contains 7 floats: (x, y, z) (w; x, y, z)
type Placement [7]float32

func (p Placement) MarshalSII() ([]byte, error) {
	var siiFloats = make([][]byte, 7)

	for i, f := range p {
		b, err := float2sii(f)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to encode float")
		}
		siiFloats[i] = b
	}

	var buf = new(bytes.Buffer)

	buf.Write([]byte("("))
	buf.Write(bytes.Join(siiFloats[0:3], []byte(", ")))
	buf.Write([]byte(") ("))
	buf.Write(siiFloats[3])
	buf.Write([]byte("; "))
	buf.Write(bytes.Join(siiFloats[4:7], []byte(", ")))
	buf.Write([]byte(")"))

	return buf.Bytes(), nil
}

func (p *Placement) UnmarshalSII(in []byte) error {
	if !placementRegexp.Match(in) {
		return errors.New("Input data does not match expected format")
	}

	grps := placementRegexp.FindSubmatch(in)
	var err error
	for i := 0; i < 7; i++ {
		if p[i], err = sii2float(grps[i+1]); err != nil {
			return errors.Wrap(err, "Unable to decode float")
		}
	}

	return nil
}

// fixed => native type int

// fixed2-4 => [2]int - [4]int

// int2 => [2]int

// quaternion => [4]float

// s16, s32, s64 => int16, int32, int64

// u16, u32, u64 => uint16, uint32, uint64

// bool => native type bool

// token => native type string

type Ptr struct {
	Target string
	unit   *Unit
}

func (p Ptr) MarshalSII() []byte { return []byte(p.Target) }

func (p Ptr) Resolve() Block {
	if p.Target == "null" {
		return nil
	}

	for _, b := range p.unit.Entries {
		if b.Name() == p.Target {
			return b
		}
	}
	return nil
}

func (p *Ptr) UnmarshalSII(in []byte) error {
	p.Target = string(in)
	return nil
}

// resource_tie => native type string

// RawValue is used in places where a key can has multiple types and
// clean parsing into Go types is no longer possible. Sadly even parsing
// into interface{} is not possible as even for that the type of the value
// must be known
type RawValue []byte

func (r RawValue) MarshalSII() ([]byte, error) { return r, nil }

func (r *RawValue) UnmarshalSII(in []byte) error {
	*r = in
	return nil
}

// TODO: Add converter functions from / to RawValue
