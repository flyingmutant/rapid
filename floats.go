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
	float32SignifBits = 23

	float64ExpBits    = 11
	float64ExpBias    = 1<<(float64ExpBits-1) - 1
	float64SignifBits = 52

	floatExpLabel    = "floatexp"
	floatSignifLabel = "floatsignif"

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
		typ:        float32Type,
		signifBits: float32SignifBits,
		maxVal:     math.MaxFloat32,
		allowInf:   allowInf,
		allowNan:   allowNan,
	})
}

func Float64sEx(allowInf bool, allowNan bool) *Generator {
	return newGenerator(&floatGen{
		typ:        float64Type,
		signifBits: float64SignifBits,
		maxVal:     math.MaxFloat64,
		allowInf:   allowInf,
		allowNan:   allowNan,
	})
}

type floatGen struct {
	typ        reflect.Type
	signifBits uint
	maxVal     float64
	allowInf   bool
	allowNan   bool
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
	f := genUfloatRange(s, 0, g.maxVal, g.signifBits)

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

func ufloatFracBits(e int32, signifBits uint) uint {
	if e <= 0 {
		return signifBits
	} else if uint(e) < signifBits {
		return signifBits - uint(e)
	} else {
		return 0
	}
}

func ufloatParts(f float64, signifBits uint) (int32, uint64, uint64) {
	u := math.Float64bits(f)
	e := int32(u>>float64SignifBits) - float64ExpBias
	b := ufloatFracBits(e, signifBits)
	s := (u & bitmask64(float64SignifBits)) >> (float64SignifBits - signifBits)
	return e, s >> b, s & bitmask64(b)
}

func genUfloatRange(s bitStream, min float64, max float64, signifBits uint) float64 {
	assert(min >= 0 && min < max)

	minExp, minSignifI, minSignifF := ufloatParts(min, signifBits)
	maxExp, maxSignifI, maxSignifF := ufloatParts(max, signifBits)

	i := s.beginGroup(floatExpLabel, false)
	e := genIntRange(s, int64(minExp), int64(maxExp), true)
	s.endGroup(i, false)

	fracBits := ufloatFracBits(int32(e), signifBits)

	j := s.beginGroup(floatSignifLabel, false)
	var siMin, siMax uint64
	switch {
	case minExp == maxExp:
		siMin, siMax = minSignifI, maxSignifI
	case int32(e) == minExp:
		siMin, siMax = minSignifI, bitmask64(signifBits-fracBits)
	case int32(e) == maxExp:
		siMin, siMax = 0, maxSignifI
	default:
		siMin, siMax = 0, bitmask64(signifBits-fracBits)
	}
	si := genUintRange(s, siMin, siMax, false)
	var sfMin, sfMax uint64
	switch {
	case minExp == maxExp && minSignifI == maxSignifI:
		sfMin, sfMax = minSignifF, maxSignifF
	case int32(e) == minExp && si == minSignifI:
		sfMin, sfMax = minSignifF, bitmask64(fracBits)
	case int32(e) == maxExp && si == maxSignifI:
		sfMin, sfMax = 0, maxSignifF
	default:
		sfMin, sfMax = 0, bitmask64(fracBits)
	}
	sf := genUintRange(s, sfMin, sfMax, true)
	s.endGroup(j, false)

	sf <<= fracBits - uint(bits.Len64(sf))
	for sf > sfMax {
		sf >>= 1
	}

	return math.Float64frombits(uint64(e+float64ExpBias)<<float64SignifBits | si<<fracBits | sf)
}
