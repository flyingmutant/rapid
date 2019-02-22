// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"math"
	"reflect"
)

const (
	float32MantBits = 23

	float64ExpBits  = 11
	float64ExpBias  = 1<<(float64ExpBits-1) - 1
	float64MantBits = 52

	floatExpLabel  = "floatexp"
	floatMantLabel = "floatmant"

	floatGenTries    = 100
	failedToGenFloat = "failed to generate suitable floating-point number"
)

var (
	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))
)

func Float32s() *Generator {
	return Float32sEx(false, false)
}

func Float64s() *Generator {
	return Float64sEx(false, false)
}

func Float32sEx(allowInf bool, allowNan bool) *Generator {
	return newGenerator(&floatGen{
		typ:      float32Type,
		mantBits: float32MantBits,
		maxVal:   math.MaxFloat32,
		allowInf: allowInf,
		allowNan: allowNan,
	})
}

func Float64sEx(allowInf bool, allowNan bool) *Generator {
	return newGenerator(&floatGen{
		typ:      float64Type,
		mantBits: float64MantBits,
		maxVal:   math.MaxFloat64,
		allowInf: allowInf,
		allowNan: allowNan,
	})
}

type floatGen struct {
	typ      reflect.Type
	mantBits uint
	maxVal   float64
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
	return satisfy(func(v Value) bool {
		f := reflect.ValueOf(v).Float()
		if !g.allowInf && (f < -g.maxVal || f > g.maxVal) {
			return false
		}
		if !g.allowNan && f != f {
			return false
		}
		return true
	}, g.value_, s, floatGenTries, failedToGenFloat)
}

func (g *floatGen) value_(s bitStream) Value {
	f := genUfloatRange(s, 0, g.maxVal, g.mantBits)

	sign := s.drawBits(1)
	if sign == 1 {
		f = -f
	}

	if g.typ == float32Type {
		return float32(f)
	} else {
		return f
	}
}

func ufloatExp(f float64) int32 {
	return int32(math.Float64bits(f)>>float64MantBits) - float64ExpBias
}

// TODO: rejection sampling is *really* bad for some ranges
func genUfloatRange(s bitStream, min float64, max float64, mantBits uint) float64 {
	assert(min >= 0 && min < max)

	i := s.beginGroup(floatExpLabel, false)
	e := genIntRange(s, int64(ufloatExp(min)), int64(ufloatExp(max)), true)
	s.endGroup(i, false)

	fracBits := uint(0)
	if e <= 0 {
		fracBits = mantBits
	} else if uint(e) < mantBits {
		fracBits = mantBits - uint(e)
	}

	for {
		i := s.beginGroup(floatMantLabel, false)
		m1 := genUintN(s, uint64(1<<uint(mantBits-fracBits)-1), false)
		m2, m2w := genUintNWidth(s, uint64(1<<uint(fracBits)-1), true)
		m := m1<<fracBits | m2<<(fracBits-uint(m2w))
		f := math.Float64frombits(uint64(e+float64ExpBias)<<float64MantBits | m)
		ok := f >= min && f <= max
		s.endGroup(i, !ok)

		if ok {
			return f
		}
	}
}
