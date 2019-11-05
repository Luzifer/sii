package sii

import "strings"

// See https://modding.scssoft.com/wiki/Documentation/Engine/Units

// string => native type string

// float => native type float

// float2-4 => [2]float - [4]float

type Placement struct {
	X, Y, Z       float64
	W, X2, Y2, Z2 float64
}

// TODO: Implement marshalling for Placement

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

func (p Ptr) CanResolve() bool   { return strings.HasPrefix(p.Target, ".") }
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
