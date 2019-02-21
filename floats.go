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
	float32ExpBits  = 8
	float32ExpBias  = 1<<(float32ExpBits-1) - 1
	float32MantBits = 23
	float32MantMask = 1<<float32MantBits - 1

	float64ExpBits  = 11
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
		mantBits: float32MantBits,
		maxVal:   math.MaxFloat32,
		allowInf: allowInf,
		allowNan: allowNan,
	})
}

func Float64sEx(allowInf bool, allowNan bool) *Generator {
	return newGenerator(&floatGen{
		typ:      float64Type,
		minExp:   -float64ExpBias,
		maxExp:   float64ExpBias + 1,
		mantBits: float64MantBits,
		maxVal:   math.MaxFloat64,
		allowInf: allowInf,
		allowNan: allowNan,
	})
}

type floatGen struct {
	typ      reflect.Type
	minExp   int32
	maxExp   int32
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
	e := genIntRange(s, int64(g.minExp), int64(g.maxExp), true)

	fracBits := uint(0)
	if e <= 0 {
		fracBits = g.mantBits
	} else if uint(e) < g.mantBits {
		fracBits = g.mantBits - uint(e)
	}

	m1 := genUintN(s, uint64(1<<uint(g.mantBits-fracBits)-1), false)
	m2, m2w := genUintNWidth(s, uint64(1<<uint(fracBits)-1), true)
	sign := s.drawBits(1)

	m := m1<<fracBits | m2<<(fracBits-uint(m2w))

	if g.typ == float32Type {
		return math.Float32frombits(uint32(sign)<<31 | uint32(e+float32ExpBias)<<float32MantBits | uint32(m))
	} else {
		return math.Float64frombits(sign<<63 | uint64(e+float64ExpBias)<<float64MantBits | m)
	}
}
