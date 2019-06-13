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
	float32ExpBits    = 8
	float32SignifBits = 23

	float64ExpBits    = 11
	float64SignifBits = 52

	floatExpLabel    = "floatexp"
	floatSignifLabel = "floatsignif"
)

var (
	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))
)

func Float32s() *Generator {
	return Float32sRange(-math.MaxFloat32, math.MaxFloat32)
}

func Float32sMin(min float32) *Generator {
	return Float32sRange(min, math.MaxFloat32)
}

func Float32sMax(max float32) *Generator {
	return Float32sRange(-math.MaxFloat32, max)
}

func Float32sRange(min float32, max float32) *Generator {
	assertf(min == min, "min should not be a NaN")
	assertf(max == max, "max should not be a NaN")
	assertf(min <= max, "invalid range [%v, %v]", min, max)

	return newGenerator(&floatGen{
		typ:        float32Type,
		expBits:    float32ExpBits,
		signifBits: float32SignifBits,
		min:        float64(min),
		max:        float64(max),
		minVal:     -math.MaxFloat32,
		maxVal:     math.MaxFloat32,
	})
}

func Float64s() *Generator {
	return Float64sRange(-math.MaxFloat64, math.MaxFloat64)
}

func Float64sMin(min float64) *Generator {
	return Float64sRange(min, math.MaxFloat64)
}

func Float64sMax(max float64) *Generator {
	return Float64sRange(-math.MaxFloat64, max)
}

func Float64sRange(min float64, max float64) *Generator {
	assertf(min == min, "min should not be a NaN")
	assertf(max == max, "max should not be a NaN")
	assertf(min <= max, "invalid range [%v, %v]", min, max)

	return newGenerator(&floatGen{
		typ:        float64Type,
		expBits:    float64ExpBits,
		signifBits: float64SignifBits,
		min:        min,
		max:        max,
		minVal:     -math.MaxFloat64,
		maxVal:     math.MaxFloat64,
	})
}

type floatGen struct {
	typ        reflect.Type
	expBits    uint
	signifBits uint
	min        float64
	max        float64
	minVal     float64
	maxVal     float64
}

func (g *floatGen) String() string {
	kind := "Float64s"
	if g.typ == float32Type {
		kind = "Float32s"
	}

	if g.min != g.minVal && g.max != g.maxVal {
		return fmt.Sprintf("%sRange(%g, %g)", kind, g.min, g.max)
	} else if g.min != g.minVal {
		return fmt.Sprintf("%sMin(%g)", kind, g.min)
	} else if g.max != g.maxVal {
		return fmt.Sprintf("%sMax(%g)", kind, g.max)
	}

	return fmt.Sprintf("%s()", kind)
}

func (g *floatGen) type_() reflect.Type {
	return g.typ
}

func (g *floatGen) value(s bitStream) Value {
	f := genFloatRange(s, g.min, g.max, g.expBits, g.signifBits)

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

func ufloatParts(f float64, expBits uint, signifBits uint) (int32, uint64, uint64) {
	u := math.Float64bits(f) & math.MaxInt64

	e := int32(u>>float64SignifBits) - int32(bitmask64(float64ExpBits-1))
	b := int32(bitmask64(expBits - 1))
	if e < -b+1 {
		e = -b + 1 // -b is subnormal
	} else if e > b {
		e = b // b+1 is Inf/NaN
	}

	s := (u & bitmask64(float64SignifBits)) >> (float64SignifBits - signifBits)
	n := ufloatFracBits(e, signifBits)

	return e, s >> n, s & bitmask64(n)
}

func ufloatFromParts(signifBits uint, e int32, si uint64, sf uint64) float64 {
	n := ufloatFracBits(e, signifBits)

	e_ := (uint64(e) + bitmask64(float64ExpBits-1)) << float64SignifBits
	s_ := (si<<n | sf) << (float64SignifBits - signifBits)

	return math.Float64frombits(e_ | s_)
}

func genUfloatRange(s bitStream, min float64, max float64, expBits uint, signifBits uint) float64 {
	assert(min >= 0 && min <= max)

	minExp, minSignifI, minSignifF := ufloatParts(min, expBits, signifBits)
	maxExp, maxSignifI, maxSignifF := ufloatParts(max, expBits, signifBits)

	i := s.beginGroup(floatExpLabel, false)
	e, lOverflow, rOverflow := genIntRange(s, int64(minExp), int64(maxExp), true)
	s.endGroup(i, false)

	fracBits := ufloatFracBits(int32(e), signifBits)

	j := s.beginGroup(floatSignifLabel, false)
	var siMin, siMax uint64
	switch {
	case lOverflow:
		siMin, siMax = minSignifI, minSignifI
	case rOverflow:
		siMin, siMax = maxSignifI, maxSignifI
	case minExp == maxExp:
		siMin, siMax = minSignifI, maxSignifI
	case int32(e) == minExp:
		siMin, siMax = minSignifI, bitmask64(signifBits-fracBits)
	case int32(e) == maxExp:
		siMin, siMax = 0, maxSignifI
	default:
		siMin, siMax = 0, bitmask64(signifBits-fracBits)
	}
	si, _, _ := genUintRange(s, siMin, siMax, false)
	var sfMin, sfMax uint64
	switch {
	case lOverflow:
		sfMin, sfMax = minSignifF, minSignifF
	case rOverflow:
		sfMin, sfMax = maxSignifF, maxSignifF
	case minExp == maxExp && minSignifI == maxSignifI:
		sfMin, sfMax = minSignifF, maxSignifF
	case int32(e) == minExp && si == minSignifI:
		sfMin, sfMax = minSignifF, bitmask64(fracBits)
	case int32(e) == maxExp && si == maxSignifI:
		sfMin, sfMax = 0, maxSignifF
	default:
		sfMin, sfMax = 0, bitmask64(fracBits)
	}
	maxR := bits.Len64(sfMax - sfMin)
	r := genUintNNoReject(s, uint64(maxR))
	sf, _, _ := genUintRange(s, sfMin, sfMax, false)
	s.endGroup(j, false)

	for i := uint(0); i < uint(maxR)-uint(r); i++ {
		mask := ^(uint64(1) << i)
		if sf&mask < sfMin {
			break
		}
		sf &= mask
	}

	return ufloatFromParts(signifBits, int32(e), si, sf)
}

func genFloatRange(s bitStream, min float64, max float64, expBits uint, signifBits uint) float64 {
	var posMin, negMin, pNeg float64
	if min >= 0 {
		posMin = min
		pNeg = 0
	} else if max <= 0 {
		negMin = -max
		pNeg = 1
	} else {
		pNeg = 0.5
	}

	if flipBiasedCoin(s, pNeg) {
		return -genUfloatRange(s, negMin, -min, expBits, signifBits)
	} else {
		return genUfloatRange(s, posMin, max, expBits, signifBits)
	}
}
