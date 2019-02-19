// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"math"
	"math/bits"
	"reflect"
)

const (
	float32ExpBits  = 8
	float32ExpMask  = 1<<float32ExpBits - 1
	float32ExpBias  = 1<<(float32ExpBits-1) - 1
	float32MantBits = 23
	float32MantMask = 1<<float32MantBits - 1

	float64ExpBits  = 11
	float64ExpMask  = 1<<float64ExpBits - 1
	float64ExpBias  = 1<<(float64ExpBits-1) - 1
	float64MantBits = 52
	float64MantMask = 1<<float64MantBits - 1

	floatGenTries    = 100
	failedToGenFloat = "failed to generate suitable floating-point number"
)

var (
	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))
)

func reverseBits(x uint64, n uint) uint64 {
	assert(n <= 64 && bits.Len64(x) <= int(n))
	return bits.Reverse64(x) >> (64 - n)
}

func transformMant(e uint32, bias uint32, bits_ uint, m uint64) uint64 {
	var fractBits uint
	if e <= bias {
		fractBits = bits_
	} else if e < bias+uint32(bits_) {
		fractBits = bits_ - uint(e-bias)
	} else {
		fractBits = 0
	}

	fractPart := m & (1<<fractBits - 1)
	m_ := (m ^ fractPart) | reverseBits(fractPart, fractBits)

	assert(bits.Len64(m_) <= int(bits_))
	return m_
}

func transformMant32(e uint32, m uint64) uint64 {
	return transformMant(e, float32ExpBias, float32MantBits, m)
}
func transformMant64(e uint32, m uint64) uint64 {
	return transformMant(e, float64ExpBias, float64MantBits, m)
}

func float32ToLex(f float32) (bool, uint32, uint64) {
	i := math.Float32bits(f)
	u := i & math.MaxInt32
	e := uint32(u >> float32MantBits)
	m := uint64(u & float32MantMask)
	return i&(1<<31) != 0, e, transformMant32(e, m)
}

func float64ToLex(f float64) (bool, uint32, uint64) {
	i := math.Float64bits(f)
	u := i & math.MaxInt64
	e := uint32(u >> float64MantBits)
	m := uint64(u & float64MantMask)
	return i&(1<<63) != 0, e, transformMant64(e, m)
}

func lexToFloat32(sign bool, e uint32, m uint64) float32 {
	e = e & float32ExpMask
	m = transformMant32(e, m&float32MantMask)
	u := e<<float32MantBits | uint32(m)
	if sign {
		u |= 1 << 31
	}
	return math.Float32frombits(u)
}

func lexToFloat64(sign bool, e uint32, m uint64) float64 {
	e = e & float64ExpMask
	m = transformMant64(e, m&float64MantMask)
	u := uint64(e)<<float64MantBits | m
	if sign {
		u |= 1 << 63
	}
	return math.Float64frombits(u)
}

func Float32s() *Generator {
	return Float32sEx(false, false)
}

func Float64s() *Generator {
	return Float64sEx(false, false)
}

func Float32sEx(allowInf bool, allowNan bool) *Generator {
	return newGenerator(&floatGen{
		typ:      float32Type,
		minExp:   -float32ExpBias,
		maxExp:   float32ExpBias + 1,
		maxMant:  float32MantMask,
		allowInf: allowInf,
		allowNan: allowNan,
	})
}

func Float64sEx(allowInf bool, allowNan bool) *Generator {
	return newGenerator(&floatGen{
		typ:      float64Type,
		minExp:   -float64ExpBias,
		maxExp:   float64ExpBias + 1,
		maxMant:  float64MantMask,
		allowInf: allowInf,
		allowNan: allowNan,
	})
}

type floatGen struct {
	typ      reflect.Type
	minExp   int32
	maxExp   int32
	maxMant  uint64
	allowInf bool
	allowNan bool
}

func (g *floatGen) String() string {
	if g.typ == float32Type {
		if !g.allowInf && !g.allowNan {
			return "Float32s()"
		} else {
			return fmt.Sprintf("Float32sEx(allowInf=%v, allowNan=%v)", g.allowInf, g.allowNan)
		}
	} else {
		if !g.allowInf && !g.allowNan {
			return "Float64s()"
		} else {
			return fmt.Sprintf("Float64sEx(allowInf=%v, allowNan=%v)", g.allowInf, g.allowNan)
		}
	}
}

func (g *floatGen) type_() reflect.Type {
	return g.typ
}

func (g *floatGen) value(s bitStream) Value {
	var cond func(Value) bool
	if g.typ == float32Type {
		cond = func(v Value) bool {
			f := v.(float32)
			if !g.allowInf && (f < -math.MaxFloat32 || f > math.MaxFloat32) {
				return false
			}
			if !g.allowNan && f != f {
				return false
			}
			return true
		}
	} else {
		cond = func(v Value) bool {
			f := v.(float64)
			if !g.allowInf && (f < -math.MaxFloat64 || f > math.MaxFloat64) {
				return false
			}
			if !g.allowNan && f != f {
				return false
			}
			return true
		}
	}

	return satisfy(cond, g.value_, s, floatGenTries, failedToGenFloat)
}

func (g *floatGen) value_(s bitStream) Value {
	var (
		e    = genIntRange(s, int64(g.minExp), int64(g.maxExp), true)
		m    = genUintN(s, g.maxMant, false)
		sign = s.drawBits(1)
	)

	if g.typ == float32Type {
		return lexToFloat32(sign == 1, uint32(float32ExpBias+e), m)
	} else {
		return lexToFloat64(sign == 1, uint32(float64ExpBias+e), m)
	}
}
